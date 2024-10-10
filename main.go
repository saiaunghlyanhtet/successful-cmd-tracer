package main
import (
	"C"
)
import (
    "bytes"
    "encoding/binary"
    "fmt"
    "log"
    "os"
    "os/signal"
    "syscall"

    "github.com/aquasecurity/libbpfgo"
)

func main() {

    // Load the BPF program
    bpfModule, err := libbpfgo.NewModuleFromFile("execve_tracer.bpf.o")
    if err != nil {
        log.Fatalf("Error loading BPF object: %v", err)
    }
    defer bpfModule.Close()

    // Load the BPF object
    if err := bpfModule.BPFLoadObject(); err != nil {
        log.Fatalf("Error loading BPF object: %v", err)
    }

    // Attach the BPF program to the kprobe
    prog, err := bpfModule.GetProgram("kprobe.sys_execve")
    if err != nil {
        log.Fatalf("Error getting program: %v", err)
    }
    
    _, err = prog.AttachKprobe("sys_execve")
    if err != nil {
        log.Fatalf("Error attaching kprobe: %v", err)
    }

    // Create a channel to receive events from the ring buffer
    eventsChannel := make(chan []byte)
    rb, err := bpfModule.InitRingBuf("events", eventsChannel)
    if err != nil {
        log.Fatalf("Error initializing ring buffer: %v", err)
    }

	rb.Start();
    // Create a channel to handle shutdown signals
    stop := make(chan os.Signal, 1)
    signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

    fmt.Println("Tracing execve calls... Ctrl-C to end.")

    // Start a goroutine to handle incoming events
    go func() {
        for {
            eventBytes := <-eventsChannel
            if len(eventBytes) < 8 { // Ensure there are enough bytes to read
                continue
            }

            // Unpack event data
            pid := int(binary.LittleEndian.Uint32(eventBytes[0:4])) // Treat first 4 bytes as LittleEndian Uint32
            comm := string(bytes.TrimRight(eventBytes[4:], "\x00")) // Remove excess 0's from comm

            fmt.Printf("PID: %d Command: %s\n", pid, comm)
        }
    }()

    // Wait for termination signal
    <-stop
    fmt.Println("Exiting...")
}

