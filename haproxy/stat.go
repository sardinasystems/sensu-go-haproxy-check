package haproxy

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/gocarina/gocsv"
)

const (
	// Frontend key
	Frontend string = "FRONTEND"
	// Backend key
	Backend string = "BACKEND"
)

// StatLine represents one line from haproxy stat report
type StatLine struct {
	// [[[cog:
	// # Header: echo "show stat" | nc -U /var/run/haproxy.sock | head -n1
	// header = '''# pxname,svname,qcur,qmax,scur,smax,slim,stot,bin,bout,dreq,dresp,ereq,econ,eresp,wretr,wredis,status,weight,act,bck,chkfail,chkdown,lastchg,downtime,qlimit,pid,iid,sid,throttle,lbtot,tracked,type,rate,rate_lim,rate_max,check_status,check_code,check_duration,hrsp_1xx,hrsp_2xx,hrsp_3xx,hrsp_4xx,hrsp_5xx,hrsp_other,hanafail,req_rate,req_rate_max,req_tot,cli_abrt,srv_abrt,comp_in,comp_out,comp_byp,comp_rsp,lastsess,last_chk,last_agt,qtime,ctime,rtime,ttime,agent_status,agent_code,agent_duration,check_desc,agent_desc,check_rise,check_fall,check_health,agent_rise,agent_fall,agent_health,addr,cookie,mode,algo,conn_rate,conn_rate_max,conn_tot,intercepted,dcon,dses,wrew,connect,reuse,cache_lookups,cache_hits,srv_icur,src_ilim,qtime_max,ctime_max,rtime_max,ttime_max,eint,idle_conn_cur,safe_conn_cur,used_conn_cur,need_conn_est,uweight,-,h2_headers_rcvd,h2_data_rcvd,h2_settings_rcvd,h2_rst_stream_rcvd,h2_goaway_rcvd,h2_detected_conn_protocol_errors,h2_detected_strm_protocol_errors,h2_rst_stream_resp,h2_goaway_resp,h2_open_connections,h2_backend_open_streams,h2_open_connections,h2_backend_open_streams,'''
	// fields = [(s.replace('# ', '').strip().title().replace('_', ''), s) for s in header.split(',') if s not in ['', '-']]
	// used = set()
	// fields = (f for f in fields if f not in used and (used.add(f) or True))
	// for field in fields:
	//     struct_name = field[0]
	//     csv_name = field[1]
	//     json_name = field[1].replace('# ', '')
	//     type_ = 'string'
	//     #if csv_name.endswith(('name', 'desc')) or csv_name in ['status', 'mode', 'check_status']:
	//     #    type_ = 'string'
	//     if csv_name in ['scur', 'slim']:	# NOTE(vermakov): we use only a few fields in check, leave others as string
	//         type_ = 'int'
	//     cog.outl(f"""{struct_name:28s} {type_:6s} `csv:"{csv_name}" json:"{json_name},omitempty"`""")
	// ]]]
	Pxname                       string `csv:"# pxname" json:"pxname,omitempty"`
	Svname                       string `csv:"svname" json:"svname,omitempty"`
	Qcur                         string `csv:"qcur" json:"qcur,omitempty"`
	Qmax                         string `csv:"qmax" json:"qmax,omitempty"`
	Scur                         int    `csv:"scur" json:"scur,omitempty"`
	Smax                         string `csv:"smax" json:"smax,omitempty"`
	Slim                         int    `csv:"slim" json:"slim,omitempty"`
	Stot                         string `csv:"stot" json:"stot,omitempty"`
	Bin                          string `csv:"bin" json:"bin,omitempty"`
	Bout                         string `csv:"bout" json:"bout,omitempty"`
	Dreq                         string `csv:"dreq" json:"dreq,omitempty"`
	Dresp                        string `csv:"dresp" json:"dresp,omitempty"`
	Ereq                         string `csv:"ereq" json:"ereq,omitempty"`
	Econ                         string `csv:"econ" json:"econ,omitempty"`
	Eresp                        string `csv:"eresp" json:"eresp,omitempty"`
	Wretr                        string `csv:"wretr" json:"wretr,omitempty"`
	Wredis                       string `csv:"wredis" json:"wredis,omitempty"`
	Status                       string `csv:"status" json:"status,omitempty"`
	Weight                       string `csv:"weight" json:"weight,omitempty"`
	Act                          string `csv:"act" json:"act,omitempty"`
	Bck                          string `csv:"bck" json:"bck,omitempty"`
	Chkfail                      string `csv:"chkfail" json:"chkfail,omitempty"`
	Chkdown                      string `csv:"chkdown" json:"chkdown,omitempty"`
	Lastchg                      string `csv:"lastchg" json:"lastchg,omitempty"`
	Downtime                     string `csv:"downtime" json:"downtime,omitempty"`
	Qlimit                       string `csv:"qlimit" json:"qlimit,omitempty"`
	Pid                          string `csv:"pid" json:"pid,omitempty"`
	Iid                          string `csv:"iid" json:"iid,omitempty"`
	Sid                          string `csv:"sid" json:"sid,omitempty"`
	Throttle                     string `csv:"throttle" json:"throttle,omitempty"`
	Lbtot                        string `csv:"lbtot" json:"lbtot,omitempty"`
	Tracked                      string `csv:"tracked" json:"tracked,omitempty"`
	Type                         string `csv:"type" json:"type,omitempty"`
	Rate                         string `csv:"rate" json:"rate,omitempty"`
	RateLim                      string `csv:"rate_lim" json:"rate_lim,omitempty"`
	RateMax                      string `csv:"rate_max" json:"rate_max,omitempty"`
	CheckStatus                  string `csv:"check_status" json:"check_status,omitempty"`
	CheckCode                    string `csv:"check_code" json:"check_code,omitempty"`
	CheckDuration                string `csv:"check_duration" json:"check_duration,omitempty"`
	Hrsp1Xx                      string `csv:"hrsp_1xx" json:"hrsp_1xx,omitempty"`
	Hrsp2Xx                      string `csv:"hrsp_2xx" json:"hrsp_2xx,omitempty"`
	Hrsp3Xx                      string `csv:"hrsp_3xx" json:"hrsp_3xx,omitempty"`
	Hrsp4Xx                      string `csv:"hrsp_4xx" json:"hrsp_4xx,omitempty"`
	Hrsp5Xx                      string `csv:"hrsp_5xx" json:"hrsp_5xx,omitempty"`
	HrspOther                    string `csv:"hrsp_other" json:"hrsp_other,omitempty"`
	Hanafail                     string `csv:"hanafail" json:"hanafail,omitempty"`
	ReqRate                      string `csv:"req_rate" json:"req_rate,omitempty"`
	ReqRateMax                   string `csv:"req_rate_max" json:"req_rate_max,omitempty"`
	ReqTot                       string `csv:"req_tot" json:"req_tot,omitempty"`
	CliAbrt                      string `csv:"cli_abrt" json:"cli_abrt,omitempty"`
	SrvAbrt                      string `csv:"srv_abrt" json:"srv_abrt,omitempty"`
	CompIn                       string `csv:"comp_in" json:"comp_in,omitempty"`
	CompOut                      string `csv:"comp_out" json:"comp_out,omitempty"`
	CompByp                      string `csv:"comp_byp" json:"comp_byp,omitempty"`
	CompRsp                      string `csv:"comp_rsp" json:"comp_rsp,omitempty"`
	Lastsess                     string `csv:"lastsess" json:"lastsess,omitempty"`
	LastChk                      string `csv:"last_chk" json:"last_chk,omitempty"`
	LastAgt                      string `csv:"last_agt" json:"last_agt,omitempty"`
	Qtime                        string `csv:"qtime" json:"qtime,omitempty"`
	Ctime                        string `csv:"ctime" json:"ctime,omitempty"`
	Rtime                        string `csv:"rtime" json:"rtime,omitempty"`
	Ttime                        string `csv:"ttime" json:"ttime,omitempty"`
	AgentStatus                  string `csv:"agent_status" json:"agent_status,omitempty"`
	AgentCode                    string `csv:"agent_code" json:"agent_code,omitempty"`
	AgentDuration                string `csv:"agent_duration" json:"agent_duration,omitempty"`
	CheckDesc                    string `csv:"check_desc" json:"check_desc,omitempty"`
	AgentDesc                    string `csv:"agent_desc" json:"agent_desc,omitempty"`
	CheckRise                    string `csv:"check_rise" json:"check_rise,omitempty"`
	CheckFall                    string `csv:"check_fall" json:"check_fall,omitempty"`
	CheckHealth                  string `csv:"check_health" json:"check_health,omitempty"`
	AgentRise                    string `csv:"agent_rise" json:"agent_rise,omitempty"`
	AgentFall                    string `csv:"agent_fall" json:"agent_fall,omitempty"`
	AgentHealth                  string `csv:"agent_health" json:"agent_health,omitempty"`
	Addr                         string `csv:"addr" json:"addr,omitempty"`
	Cookie                       string `csv:"cookie" json:"cookie,omitempty"`
	Mode                         string `csv:"mode" json:"mode,omitempty"`
	Algo                         string `csv:"algo" json:"algo,omitempty"`
	ConnRate                     string `csv:"conn_rate" json:"conn_rate,omitempty"`
	ConnRateMax                  string `csv:"conn_rate_max" json:"conn_rate_max,omitempty"`
	ConnTot                      string `csv:"conn_tot" json:"conn_tot,omitempty"`
	Intercepted                  string `csv:"intercepted" json:"intercepted,omitempty"`
	Dcon                         string `csv:"dcon" json:"dcon,omitempty"`
	Dses                         string `csv:"dses" json:"dses,omitempty"`
	Wrew                         string `csv:"wrew" json:"wrew,omitempty"`
	Connect                      string `csv:"connect" json:"connect,omitempty"`
	Reuse                        string `csv:"reuse" json:"reuse,omitempty"`
	CacheLookups                 string `csv:"cache_lookups" json:"cache_lookups,omitempty"`
	CacheHits                    string `csv:"cache_hits" json:"cache_hits,omitempty"`
	SrvIcur                      string `csv:"srv_icur" json:"srv_icur,omitempty"`
	SrcIlim                      string `csv:"src_ilim" json:"src_ilim,omitempty"`
	QtimeMax                     string `csv:"qtime_max" json:"qtime_max,omitempty"`
	CtimeMax                     string `csv:"ctime_max" json:"ctime_max,omitempty"`
	RtimeMax                     string `csv:"rtime_max" json:"rtime_max,omitempty"`
	TtimeMax                     string `csv:"ttime_max" json:"ttime_max,omitempty"`
	Eint                         string `csv:"eint" json:"eint,omitempty"`
	IdleConnCur                  string `csv:"idle_conn_cur" json:"idle_conn_cur,omitempty"`
	SafeConnCur                  string `csv:"safe_conn_cur" json:"safe_conn_cur,omitempty"`
	UsedConnCur                  string `csv:"used_conn_cur" json:"used_conn_cur,omitempty"`
	NeedConnEst                  string `csv:"need_conn_est" json:"need_conn_est,omitempty"`
	Uweight                      string `csv:"uweight" json:"uweight,omitempty"`
	H2HeadersRcvd                string `csv:"h2_headers_rcvd" json:"h2_headers_rcvd,omitempty"`
	H2DataRcvd                   string `csv:"h2_data_rcvd" json:"h2_data_rcvd,omitempty"`
	H2SettingsRcvd               string `csv:"h2_settings_rcvd" json:"h2_settings_rcvd,omitempty"`
	H2RstStreamRcvd              string `csv:"h2_rst_stream_rcvd" json:"h2_rst_stream_rcvd,omitempty"`
	H2GoawayRcvd                 string `csv:"h2_goaway_rcvd" json:"h2_goaway_rcvd,omitempty"`
	H2DetectedConnProtocolErrors string `csv:"h2_detected_conn_protocol_errors" json:"h2_detected_conn_protocol_errors,omitempty"`
	H2DetectedStrmProtocolErrors string `csv:"h2_detected_strm_protocol_errors" json:"h2_detected_strm_protocol_errors,omitempty"`
	H2RstStreamResp              string `csv:"h2_rst_stream_resp" json:"h2_rst_stream_resp,omitempty"`
	H2GoawayResp                 string `csv:"h2_goaway_resp" json:"h2_goaway_resp,omitempty"`
	H2OpenConnections            string `csv:"h2_open_connections" json:"h2_open_connections,omitempty"`
	H2BackendOpenStreams         string `csv:"h2_backend_open_streams" json:"h2_backend_open_streams,omitempty"`
	// [[[end]]] (checksum: 0317562b38ec931796a67e29935fb72f)
}

// StatService is a mapping of all SvName lines
type StatService map[string]StatLine

// Stats is a mapping for PxName -> SvName -> StatLine
type Stats = map[string]StatService

// ParseStatCSV parses stats csv into Stats
func ParseStatCSV(data io.Reader) (Stats, []byte, error) {
	lines := []StatLine{}

	rawData, err := io.ReadAll(data)
	if err != nil {
		return nil, nil, fmt.Errorf("read error: %w", err)
	}

	err = gocsv.UnmarshalBytes(rawData, &lines)
	if err != nil {
		return nil, nil, fmt.Errorf("csv parse error: %w", err)
	}

	out := make(Stats)
	for _, line := range lines {
		pxmap, ok := out[line.Pxname]
		if !ok {
			pxmap = make(StatService)
			out[line.Pxname] = pxmap
		}

		pxmap[line.Svname] = line
	}

	return out, rawData, nil
}

// GetStats query HAProxy for Stats
func GetStats(socketPath string) (Stats, []byte, error) {
	addr := &net.UnixAddr{Name: socketPath}
	sock, err := net.DialUnix("unix", nil, addr)
	if err != nil {
		return nil, nil, fmt.Errorf("socket open error: %w", err)
	}
	defer sock.Close()

	// Suggest that IO shouldn't ever reach so long timeout
	err = sock.SetDeadline(time.Now().Add(time.Second))
	if err != nil {
		return nil, nil, fmt.Errorf("socket set deadline error: %w", err)
	}

	_, err = sock.Write([]byte("show stat\n"))
	if err != nil {
		return nil, nil, fmt.Errorf("socket request error: %w", err)
	}

	return ParseStatCSV(sock)
}

// IsUp checks that status of the service is up
func (l StatLine) IsUp(backend *StatLine) bool {
	// In some rare calls we get empty Status for servers
	// in that case we fallback to backend state
	if l.Status == "" && backend != nil {
		return backend.IsUp(nil)
	}

	// XXX FIXME(vermakov): revise that later, observed on HAproxy 2.3.0 -- 2.3.2
	// sometimes we got report without BACKEND and empty Status
	// let's consider it as ok, otherwise we get very noisy false positive notification
	if l.Status == "" && backend == nil {
		return true
	}

	return (l.Status == "OPEN" ||
		l.Status == "UP" ||
		l.Status == "no check" ||
		l.Status == "DRAIN")
}

// LogName make a name for check logs
func (l StatLine) LogName() string {
	if l.CheckStatus == "" {
		return fmt.Sprintf("%s/%s", l.Pxname, l.Svname)
	}

	return fmt.Sprintf("%s/%s[%s]", l.Pxname, l.Svname, l.CheckStatus)
}

// SessionLimitPercentage calculates percentage usage of sessions limit
func (l StatLine) SessionLimitPercentage() float32 {
	return 100.0 * float32(l.Scur) / float32(l.Slim)
}

// Servers makes a copy of StatService without frontend and backend entries
func (s StatService) Servers() StatService {
	return s.Filter(func(s StatLine) bool {
		// XXX(vermakov): we also filter empty Svname because that must be an error in HAproxy 2.3.0+
		return s.Svname != Frontend && s.Svname != Backend && s.Svname != ""
	})
}

// Filter return entries which passes the testFunc
func (s StatService) Filter(testFunc func(StatLine) bool) StatService {
	ret := make(StatService)
	for key, value := range s {
		if !testFunc(value) {
			continue
		}

		ret[key] = value
	}

	return ret
}
