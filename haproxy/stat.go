package haproxy

import (
	"fmt"
	"io"

	"github.com/gocarina/gocsv"
)

type StatLine struct {
	// [[[cog:
	// header = '''# pxname,svname,qcur,qmax,scur,smax,slim,stot,bin,bout,dreq,dresp,ereq,econ,eresp,wretr,wredis,status,weight,act,bck,chkfail,chkdown,lastchg,downtime,qlimit,pid,iid,sid,throttle,lbtot,tracked,type,rate,rate_lim,rate_max,check_status,check_code,check_duration,hrsp_1xx,hrsp_2xx,hrsp_3xx,hrsp_4xx,hrsp_5xx,hrsp_other,hanafail,req_rate,req_rate_max,req_tot,cli_abrt,srv_abrt,comp_in,comp_out,comp_byp,comp_rsp,lastsess,last_chk,last_agt,qtime,ctime,rtime,ttime,agent_status,agent_code,agent_duration,check_desc,agent_desc,check_rise,check_fall,check_health,agent_rise,agent_fall,agent_health,addr,cookie,mode,algo,conn_rate,conn_rate_max,conn_tot,intercepted,dcon,dses,wrew,connect,reuse,cache_lookups,cache_hits,srv_icur,src_ilim,qtime_max,ctime_max,rtime_max,ttime_max,eint,idle_conn_cur,safe_conn_cur,used_conn_cur,need_conn_est,uweight,-,h2_headers_rcvd,h2_data_rcvd,h2_settings_rcvd,h2_rst_stream_rcvd,h2_goaway_rcvd,h2_detected_conn_protocol_errors,h2_detected_strm_protocol_errors,h2_rst_stream_resp,h2_goaway_resp,h2_open_connections,h2_backend_open_streams,h2_open_connections,h2_backend_open_streams,'''
	// fields = [(s.replace('# ', '').strip().title().replace('_', ''), s) for s in header.split(',') if s not in ['', '-']]
	// used = set()
	// fields = (f for f in fields if f not in used and (used.add(f) or True))
	// for field in fields:
	//     cog.outl(f"""{field[0]} string `csv:"{field[1]}"`""")
	// ]]]
	Pxname                       string `csv:"# pxname"`
	Svname                       string `csv:"svname"`
	Qcur                         string `csv:"qcur"`
	Qmax                         string `csv:"qmax"`
	Scur                         string `csv:"scur"`
	Smax                         string `csv:"smax"`
	Slim                         string `csv:"slim"`
	Stot                         string `csv:"stot"`
	Bin                          string `csv:"bin"`
	Bout                         string `csv:"bout"`
	Dreq                         string `csv:"dreq"`
	Dresp                        string `csv:"dresp"`
	Ereq                         string `csv:"ereq"`
	Econ                         string `csv:"econ"`
	Eresp                        string `csv:"eresp"`
	Wretr                        string `csv:"wretr"`
	Wredis                       string `csv:"wredis"`
	Status                       string `csv:"status"`
	Weight                       string `csv:"weight"`
	Act                          string `csv:"act"`
	Bck                          string `csv:"bck"`
	Chkfail                      string `csv:"chkfail"`
	Chkdown                      string `csv:"chkdown"`
	Lastchg                      string `csv:"lastchg"`
	Downtime                     string `csv:"downtime"`
	Qlimit                       string `csv:"qlimit"`
	Pid                          string `csv:"pid"`
	Iid                          string `csv:"iid"`
	Sid                          string `csv:"sid"`
	Throttle                     string `csv:"throttle"`
	Lbtot                        string `csv:"lbtot"`
	Tracked                      string `csv:"tracked"`
	Type                         string `csv:"type"`
	Rate                         string `csv:"rate"`
	RateLim                      string `csv:"rate_lim"`
	RateMax                      string `csv:"rate_max"`
	CheckStatus                  string `csv:"check_status"`
	CheckCode                    string `csv:"check_code"`
	CheckDuration                string `csv:"check_duration"`
	Hrsp1Xx                      string `csv:"hrsp_1xx"`
	Hrsp2Xx                      string `csv:"hrsp_2xx"`
	Hrsp3Xx                      string `csv:"hrsp_3xx"`
	Hrsp4Xx                      string `csv:"hrsp_4xx"`
	Hrsp5Xx                      string `csv:"hrsp_5xx"`
	HrspOther                    string `csv:"hrsp_other"`
	Hanafail                     string `csv:"hanafail"`
	ReqRate                      string `csv:"req_rate"`
	ReqRateMax                   string `csv:"req_rate_max"`
	ReqTot                       string `csv:"req_tot"`
	CliAbrt                      string `csv:"cli_abrt"`
	SrvAbrt                      string `csv:"srv_abrt"`
	CompIn                       string `csv:"comp_in"`
	CompOut                      string `csv:"comp_out"`
	CompByp                      string `csv:"comp_byp"`
	CompRsp                      string `csv:"comp_rsp"`
	Lastsess                     string `csv:"lastsess"`
	LastChk                      string `csv:"last_chk"`
	LastAgt                      string `csv:"last_agt"`
	Qtime                        string `csv:"qtime"`
	Ctime                        string `csv:"ctime"`
	Rtime                        string `csv:"rtime"`
	Ttime                        string `csv:"ttime"`
	AgentStatus                  string `csv:"agent_status"`
	AgentCode                    string `csv:"agent_code"`
	AgentDuration                string `csv:"agent_duration"`
	CheckDesc                    string `csv:"check_desc"`
	AgentDesc                    string `csv:"agent_desc"`
	CheckRise                    string `csv:"check_rise"`
	CheckFall                    string `csv:"check_fall"`
	CheckHealth                  string `csv:"check_health"`
	AgentRise                    string `csv:"agent_rise"`
	AgentFall                    string `csv:"agent_fall"`
	AgentHealth                  string `csv:"agent_health"`
	Addr                         string `csv:"addr"`
	Cookie                       string `csv:"cookie"`
	Mode                         string `csv:"mode"`
	Algo                         string `csv:"algo"`
	ConnRate                     string `csv:"conn_rate"`
	ConnRateMax                  string `csv:"conn_rate_max"`
	ConnTot                      string `csv:"conn_tot"`
	Intercepted                  string `csv:"intercepted"`
	Dcon                         string `csv:"dcon"`
	Dses                         string `csv:"dses"`
	Wrew                         string `csv:"wrew"`
	Connect                      string `csv:"connect"`
	Reuse                        string `csv:"reuse"`
	CacheLookups                 string `csv:"cache_lookups"`
	CacheHits                    string `csv:"cache_hits"`
	SrvIcur                      string `csv:"srv_icur"`
	SrcIlim                      string `csv:"src_ilim"`
	QtimeMax                     string `csv:"qtime_max"`
	CtimeMax                     string `csv:"ctime_max"`
	RtimeMax                     string `csv:"rtime_max"`
	TtimeMax                     string `csv:"ttime_max"`
	Eint                         string `csv:"eint"`
	IdleConnCur                  string `csv:"idle_conn_cur"`
	SafeConnCur                  string `csv:"safe_conn_cur"`
	UsedConnCur                  string `csv:"used_conn_cur"`
	NeedConnEst                  string `csv:"need_conn_est"`
	Uweight                      string `csv:"uweight"`
	H2HeadersRcvd                string `csv:"h2_headers_rcvd"`
	H2DataRcvd                   string `csv:"h2_data_rcvd"`
	H2SettingsRcvd               string `csv:"h2_settings_rcvd"`
	H2RstStreamRcvd              string `csv:"h2_rst_stream_rcvd"`
	H2GoawayRcvd                 string `csv:"h2_goaway_rcvd"`
	H2DetectedConnProtocolErrors string `csv:"h2_detected_conn_protocol_errors"`
	H2DetectedStrmProtocolErrors string `csv:"h2_detected_strm_protocol_errors"`
	H2RstStreamResp              string `csv:"h2_rst_stream_resp"`
	H2GoawayResp                 string `csv:"h2_goaway_resp"`
	H2OpenConnections            string `csv:"h2_open_connections"`
	H2BackendOpenStreams         string `csv:"h2_backend_open_streams"`
	// [[[end]]]
}

type StatData = map[string]map[string]StatLine

func ParseStatCSV(data io.Reader) (StatData, error) {
	lines := []StatLine{}

	err := gocsv.Unmarshal(data, &lines)
	if err != nil {
		return nil, fmt.Errorf("csv parse error: %w", err)
	}

	out := make(StatData)
	for _, line := range lines {
		pxmap, ok := out[line.Pxname]
		if !ok {
			pxmap = make(map[string]StatLine)
			out[line.Pxname] = pxmap
		}

		pxmap[line.Svname] = line
	}

	return out, nil
}
