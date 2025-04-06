
## ✅ **System Goal**
Simulate an **IoT Sensor System** that:
- Generates real-time sensor data.
- Sends data to a gRPC server using multiple concurrent workers.
- Handles responses in a structured and controlled way.

---

## 🧠 **Architecture Overview**

```txt
[Producer] --> [inputChan] --> [Processor Workers] --> [gRPC Server]
                                                        ↓
                                                  [outputChan]
                                                        ↓
                                               [Response Consumer]
```

---

## 🧩 **Core Components & Theory**

### 1. **Data Producer**
- Continuously generates sensor data (UUID, temp, timestamp).
- Uses `inputChan` to pass data to processors.
- Runs for a **limited time** (e.g. 1 minute) using a `done` channel.

### 2. **inputChan**
- Unbuffered or lightly buffered channel (`chan *proto.SensorData`).
- Acts as a **queue** between producer and processor.
- Keeps producers & processors in sync.

---

### 3. **Concurrency Control – `inputCountChan`**
- Buffered channel of type `chan struct{}`.
- Acts as a **semaphore** (e.g. size 5) to **limit concurrent gRPC calls**.
- Prevents overloading server with 100s of goroutines.

```go
inputCountChan <- struct{}{} // Acquires slot
<- inputCountChan            // Releases slot
```

---

### 4. **Processor Workers**
- Goroutine that:
  - Reads data from `inputChan`.
  - Sends it to the gRPC server.
  - Puts response into `outputChan`.
- Uses `context.WithTimeout()` to limit request time.
- Each gRPC call is protected by the semaphore.

---

### 5. **outputChan**
- Unbuffered `chan *proto.ServerResponse` (fine as long as response is fast).
- Collects gRPC responses.

---

### 6. **Response Consumer**
- Loops over `outputChan` using `for range`.
- Logs or processes responses.
- Terminates cleanly when `outputChan` is closed.

---

## 🛑 **Termination Logic**
- `Producer` stops after 1 minute (`time.Since(start)`).
- Then closes `inputChan`.
- `Processor` exits when `inputChan` is closed and drained.
- Once all processors finish, **main** closes `outputChan`.
- `Response Consumer` then finishes its loop and exits cleanly.

---

## 🔁 **Channel Behavior**
- **Unbuffered channel** blocks until receiver is ready.
- **Buffered channel** lets producer push ahead until full.
- Closing a channel signals consumers, and `range` loop ends once fully drained.

---

## 🌟 Summary in 1 Line
> "A clean, concurrent pipeline using Go channels, semaphores, and time-bound control to simulate real-time gRPC communication — with graceful startup, concurrency, and shutdown."

---
