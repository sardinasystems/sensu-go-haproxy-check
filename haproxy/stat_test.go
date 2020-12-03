package haproxy

import (
	"context"
	"net"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const testingCSV = `# pxname,svname,qcur,qmax,scur,smax,slim,stot,bin,bout,dreq,dresp,ereq,econ,eresp,wretr,wredis,status,weight,act,bck,chkfail,chkdown,lastchg,downtime,qlimit,pid,iid,sid,throttle,lbtot,tracked,type,rate,rate_lim,rate_max,check_status,check_code,check_duration,hrsp_1xx,hrsp_2xx,hrsp_3xx,hrsp_4xx,hrsp_5xx,hrsp_other,hanafail,req_rate,req_rate_max,req_tot,cli_abrt,srv_abrt,comp_in,comp_out,comp_byp,comp_rsp,lastsess,last_chk,last_agt,qtime,ctime,rtime,ttime,agent_status,agent_code,agent_duration,check_desc,agent_desc,check_rise,check_fall,check_health,agent_rise,agent_fall,agent_health,addr,cookie,mode,algo,conn_rate,conn_rate_max,conn_tot,intercepted,dcon,dses,wrew,connect,reuse,cache_lookups,cache_hits,srv_icur,src_ilim,qtime_max,ctime_max,rtime_max,ttime_max,eint,idle_conn_cur,safe_conn_cur,used_conn_cur,need_conn_est,uweight,-,h2_headers_rcvd,h2_data_rcvd,h2_settings_rcvd,h2_rst_stream_rcvd,h2_goaway_rcvd,h2_detected_conn_protocol_errors,h2_detected_strm_protocol_errors,h2_rst_stream_resp,h2_goaway_resp,h2_open_connections,h2_backend_open_streams,h2_open_connections,h2_backend_open_streams,
http,FRONTEND,,,0,0,100000,0,0,0,0,0,0,,,,,OPEN,,,,,,,,,1,2,0,,,,0,0,0,0,,,,0,0,0,0,0,0,,0,0,0,,,0,0,0,0,,,,,,,,,,,,,,,,,,,,,http,,0,0,0,0,0,0,0,,,0,0,,,,,,,0,,,,,,-,0,0,0,0,0,0,0,0,0,0,0,0,0,
https,FRONTEND,,,0,23,100000,3193,2094418,94972540,0,0,0,,,,,OPEN,,,,,,,,,1,3,0,,,,0,0,0,7,,,,0,3406,113,12,0,0,,0,7,3531,,,0,0,0,0,,,,,,,,,,,,,,,,,,,,,http,,0,7,3193,0,0,0,0,,,0,0,,,,,,,0,,,,,,-,184,328,3,0,0,0,0,0,1,0,0,1,184,
bk_dashboard_cluster,ctrl01,0,0,0,0,,0,0,0,,0,,0,0,0,0,UP,1,1,0,0,0,11225,0,,1,57,1,,0,,2,0,,0,* L7OK,0,158,0,0,0,0,0,0,,,,0,0,0,,,,,-1,,,0,0,0,0,,,,Layer7 check passed,,2,5,6,,,,,,http,,,,,,,,0,0,0,,,0,,0,0,0,0,0,0,0,0,1,1,-,,,,,,,,,,,,,,
bk_dashboard_cluster,ctrl02,0,0,0,0,,0,0,0,,0,,0,0,0,0,UP,1,1,0,0,0,11225,0,,1,57,2,,0,,2,0,,0,L7OK,200,197,0,0,0,0,0,0,,,,0,0,0,,,,,-1,,,0,0,0,0,,,,Layer7 check passed,,2,5,6,,,,,,http,,,,,,,,0,0,0,,,0,,0,0,0,0,0,0,0,0,1,1,-,,,,,,,,,,,,,,
bk_dashboard_cluster,ctrl03,0,0,0,0,,0,0,0,,0,,0,0,0,0,UP,1,1,0,0,0,11225,0,,1,57,3,,0,,2,0,,0,L7OK,200,124,0,0,0,0,0,0,,,,0,0,0,,,,,-1,,,0,0,0,0,,,,Layer7 check passed,,2,5,6,,,,,,http,,,,,,,,0,0,0,,,0,,0,0,0,0,0,0,0,0,1,1,-,,,,,,,,,,,,,,
bk_dashboard_cluster,BACKEND,0,0,0,0,10000,0,0,0,0,0,,0,0,0,0,UP,3,3,0,,0,11225,0,,1,57,0,,0,,1,0,,0,,,,0,0,0,0,0,0,,,,0,0,0,0,0,0,0,-1,,,0,0,0,0,,,,,,,,,,,,,,http,,,,,,,,0,0,0,0,0,,,0,0,0,0,0,,,,,3,-,0,0,0,0,0,0,0,0,0,0,0,0,0,
ipmi_exporter,FRONTEND,,,6,7,100000,6,1337424,6430583,0,0,0,,,,,OPEN,,,,,,,,,1,64,0,,,,0,0,0,3,,,,0,5049,0,0,0,0,,0,3,5049,,,0,0,0,0,,,,,,,,,,,,,,,,,,,,,http,,0,3,6,0,0,0,0,,,0,0,,,,,,,0,,,,,,-,0,0,0,0,0,0,0,0,0,0,0,0,0,
ipmi_exporter,ctrl01,0,0,0,2,,1683,445808,2145275,,0,,0,0,0,0,UP,1,1,0,0,0,11225,0,,1,64,1,,1683,,2,0,,1,L7OK,200,1,0,1683,0,0,0,0,,,,1683,0,0,,,,,5,,,0,0,1752,1752,,,,Layer7 check passed,,2,5,6,,,,,,http,,,,,,,,0,1019,664,,,1,,0,0,2981,2981,0,1,0,0,2,1,-,,,,,,,,,,,,,,
ipmi_exporter,ctrl02,0,0,0,2,,1683,445808,2145301,,0,,0,0,0,0,UP,1,1,0,0,0,11225,0,,1,64,2,,1683,,2,0,,1,L7OK,200,0,0,1683,0,0,0,0,,,,1683,0,0,,,,,5,,,0,0,1742,1742,,,,Layer7 check passed,,2,5,6,,,,,,http,,,,,,,,0,1018,665,,,1,,0,0,3300,3300,0,1,0,0,2,1,-,,,,,,,,,,,,,,
ipmi_exporter,ctrl03,0,0,0,2,,1683,445808,2140007,,0,,0,0,0,0,UP,1,1,0,0,0,11225,0,,1,64,3,,1683,,2,0,,1,L7OK,200,0,0,1683,0,0,0,0,,,,1683,0,0,,,,,5,,,0,0,1642,1642,,,,Layer7 check passed,,2,5,6,,,,,,http,,,,,,,,0,1001,682,,,0,,0,3,4992,4992,0,0,0,0,2,1,-,,,,,,,,,,,,,,
ipmi_exporter,BACKEND,0,0,0,6,10000,5049,1337424,6430583,0,0,,0,0,0,0,UP,3,3,0,,0,11225,0,,1,64,0,,5049,,1,0,,3,,,,0,5049,0,0,0,0,,,,5049,0,0,0,0,0,0,5,,,0,0,1711,1711,,,,,,,,,,,,,,http,,,,,,,,0,3038,2011,0,0,,,0,3,4992,4992,0,,,,,3,-,0,0,0,0,0,0,0,0,0,0,0,0,0,
`

func TestParseStatCSV(t *testing.T) {
	assert := assert.New(t)

	csvData := strings.NewReader(testingCSV)

	stats, _, err := ParseStatCSV(csvData)
	assert.NoError(err)
	assert.Len(stats, 4)

	//t.Log(stats)
}

func TestGetStats(t *testing.T) {
	assert := assert.New(t)
	tempDir := t.TempDir()

	socketPath := filepath.Join(tempDir, "haproxy.sock")

	serverCtx, serverCf := context.WithCancel(context.TODO())
	defer serverCf()

	serverReady := make(chan bool, 1)

	go func() {
		addr := &net.UnixAddr{Name: socketPath}
		ln, err := net.ListenUnix("unix", addr)
		assert.NoError(err)
		defer ln.Close()
		serverReady <- true

		for {
			select {
			case <-serverCtx.Done():
				t.Log("socket server terminated")
				return

			default:
				if err := ln.SetDeadline(time.Now().Add(time.Second)); err != nil {
					t.Fatal(err)
					return
				}

				fd, err := ln.Accept()
				if err != nil {
					if os.IsTimeout(err) {
						continue
					}

					t.Fatal(err)
				}

				go func(c net.Conn) {
					defer c.Close()

					buf := make([]byte, 1024)

					nr, err := c.Read(buf)
					assert.NoError(err)

					cmd := string(buf[0:nr])
					assert.Equal("show stat\n", cmd)

					_, err = c.Write([]byte(testingCSV))
					assert.NoError(err)
				}(fd)
			}
		}
	}()

	<-serverReady

	stats, _, err := GetStats(socketPath)
	assert.NoError(err)
	assert.Len(stats, 4)

	//t.Log(stats)
}
