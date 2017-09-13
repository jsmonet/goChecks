package main

import (
   "fmt"
   "syscall"
   "flag"
   "os"
)

type DiskStatus struct {
   All  uint64 `json:"all"`
   Free uint64 `json:"free"`
   Used uint64 `json:"used"`
 }

// disk usage of path/disk
func DiskUsage(path string) (disk DiskStatus) {
   fs := syscall.Statfs_t{}
   err := syscall.Statfs(path, &fs)
   if err != nil {
      return
   }
   disk.All = fs.Blocks * uint64(fs.Bsize)
   disk.Free = fs.Bfree * uint64(fs.Bsize)
   disk.Used = disk.All - disk.Free
   return
}

func main() {

  volumeArg := flag.String("volume", "/", "a FORMIDABLE string")
  sizeWarn := flag.Float64("warn", 75, "Enter a percentage, without the percent sign, of disk use below which you want the output to be WARN")
  sizeCrit := flag.Float64("crit", 80, "Enter a percentage, without the percent sign, of disk use below which you want the output to be CRIT")
  flag.Parse()
  disk := DiskUsage(*volumeArg)

  percentUsed := float64(disk.Used)/float64(disk.All)*100
  perUsedString := fmt.Sprintf("%.2f%% used", percentUsed)
  if percentUsed < *sizeWarn {
    fmt.Println("OK -", perUsedString)
    os.Exit(0)
  } else if percentUsed > *sizeWarn && percentUsed < *sizeCrit {
    fmt.Println("Warn -", perUsedString)
    os.Exit(1)
  } else if percentUsed > *sizeCrit {
    fmt.Println("Critical -", perUsedString)
    os.Exit(2)
  }
}
