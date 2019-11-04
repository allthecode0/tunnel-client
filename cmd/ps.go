package cmd

import (
  "github.com/hako/durafmt"
  "github.com/jedib0t/go-pretty/table"
  "github.com/labstack/tunnel-client/daemon"
  "github.com/spf13/cobra"
  "os"
  "time"
)

var psCmd = &cobra.Command{
  Use:   "ps",
  Short: "List all the connections",
  Run: func(cmd *cobra.Command, args []string) {
    psRPC()
  },
}

func psRPC() {
  startDaemon()
  c, err := getClient()
  if err != nil {
    exit(err)
  }
  defer c.Close()
  req := new(daemon.PSRequest)
  rep := new(daemon.PSReply)
  s.Start()
  err = c.Call("Server.PS", req, rep)
  if err != nil {
    exit(err)
  }
  s.Stop()
  tbl := table.NewWriter()
  tbl.SetOutputMirror(os.Stdout)
  tbl.AppendHeader(table.Row{"Name", "Target Address", "Remote URI", "Status", "Uptime"})
  for _, c := range rep.Connections {
    uptime := "-"
    if c.Status == daemon.ConnectionStatusStatusOnline {
      uptime = durafmt.ParseShort(time.Since(c.UpdatedAt)).String()
    }
    tbl.AppendRow([]interface{}{c.Name, c.TargetAddress, c.RemoteURI, c.Status, uptime})
  }
  tbl.Render()
}

func init() {
  rootCmd.AddCommand(psCmd)
}
