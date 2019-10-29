package cmangos

import (
  "fmt"
  "net"
  "time"

  "github.com/CedricThomas/api-cmangos/modules/logger"
)

type DaemonInfo struct {
  Address string `json:"address,omitempty"`
  Port int `json:"port,omitempty"`
  State int `json:"state"`
  Lastcheck int `json:"lastcheck"`
}

func CheckDaemon(d *DaemonInfo, timeout time.Duration) error {
  logger.Info(fmt.Sprintf("Check daemon state %s:%d", d.Address, d.Port))
  d.Lastcheck = int(time.Now().Unix())
  dialer := net.Dialer{Timeout: timeout * time.Second}
  c, err := dialer.Dial("tcp", fmt.Sprintf("%s:%d", d.Address, d.Port))
  if err != nil {
    logger.Error(fmt.Sprintf("Cannot connect to daemon %s:%d",
      d.Address, d.Port))
    logger.Debug(fmt.Sprintf("%v", err))
    d.State = 0
    return err
  }
  c.Close()

  d.State = 1
  return nil
}
