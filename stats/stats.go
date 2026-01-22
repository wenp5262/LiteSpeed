package stats

import "net"

// In the trimmed build, statistics collection is not required.
// These helpers keep the original call sites intact while behaving as pass-through wrappers.

func NewConn(c net.Conn) net.Conn { return c }

func NewStatsConn(c net.Conn) net.Conn { return c }

func NewStatsPacketConn(pc net.PacketConn) net.PacketConn { return pc }
