package data

import(
        "fmt"
        "math/rand"
        "sort"
        "time"
        "github.com/golang/glog"
)

type Packets []*Packet
/*func (ps Packets) Len() int {
        return len(ps)
}
func (ps Packets) Swap(i, j int) {
        ps[i], ps[j] = ps[j], ps[i]
}
func (ps Packets) Less(i,j int) bool {
        return ps[i].Timestamp.Before(ps[j].Timestamp)  // latest first
        //return ps[i].Timestamp.After(ps[j].Timestamp)  // earliest first
}*/

// Does not check for zero length slices
func (ps Packets) Summarise(table string) error /**Summary*/ {
        voltage := func(p1, p2 *Packet) bool {
                return p1.Voltage < p2.Voltage // increasing voltage
        }
        freq := func(p1, p2 *Packet) bool {
                return p1.Frequency < p2.Frequency // increasing frequency
        }
        pwr := func(p1, p2 *Packet) bool {
                return p1.ActiPwr < p2.ActiPwr // increasing power
        }
        timestamp := func(p1, p2 *Packet) bool {
                return p1.Timestamp.After(p2.Timestamp) // latest first
        }
        OrderedBy(voltage).Sort(ps)
        minv := ps[0].Voltage
        maxv := ps[len(ps)-1].Voltage
        OrderedBy(freq).Sort(ps)
        minf := ps[0].Frequency
        maxf := ps[len(ps)-1].Frequency
        OrderedBy(pwr).Sort(ps)
        minp := ps[0].ActiPwr
        maxp := ps[len(ps)-1].ActiPwr
        OrderedBy(timestamp).Sort(ps)
        summary := &Summary{Hash: ps[0].Id, Timestamp: ps[0].Timestamp, VoltageMax: maxv, VoltageMin: minv, VoltageAve: 0.0, VoltageExcpt: 0, FrequencyMax: maxf, FrequencyMin: minf, FrequencyAve: 0.0, FrequencyExcpt: 0, ActiPwrMax: maxp, ActiPwrMin: minp, ActiPwrAve: 0.0, ImportActiveEnergy: ps[0].ImActEn , ExportActiveEnergy: ps[0].ExActEn, ImportReactiveEnergy: ps[0].ImRctEn, ExportReactiveEnergy: ps[0].ExRctEn, TotalActiveEnergy: ps[0].TlActEn, TotalReactiveEnergy: ps[0].TlRctEn}
        _, err := PutSummary(table, summary)
        if err != nil {
                glog.Errorf("%v \n", err)
                return err
        }
        return nil
}

//PutPackets is a utility function to populate a packet table with samples packets for test purposes
func PutPackets(n int) error {
        db, err := GetDB()
        if err != nil {
                glog.Error(err)
                return err
        }
        start := time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC)
        for j:=13579; j<13589; j++ {
                ip := start
                for i:=0; i<n; i++ {
                        ip = ip.Add(time.Minute*5)
                        packet := &Packet{Id:int64(j), Timestamp: ip, Voltage: 240.0+rand.Float64(), Frequency: 49.0+rand.Float64(), Status: false, ActiPwr:100*rand.Float64(), AppaPwr:50*rand.Float64(), ReacPwr: 50.0*rand.Float64(), PwrFctr: rand.Float64(), ImActEn: 2.0+rand.Float64(), ExActEn:2.0*rand.Float64(), ImRctEn: 2.0*rand.Float64(), ExRctEn:2.0*rand.Float64(), TlActEn: 2.0*rand.Float64(), TlRctEn: 2.0*rand.Float64()}
                        result, err := db.Exec("insert into packet (pub_hash, created_at, voltage, frequency, protected, active_power, apparent_power, reactive_power, power_factor, import_active_energy, export_active_energy, import_reactive_energy, export_reactive_energy, total_active_energy, total_reactive_energy) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)", packet.Id, packet.Timestamp, packet.Voltage, packet.Frequency, packet.Status, packet.ActiPwr, packet.AppaPwr, packet.ReacPwr, packet.PwrFctr, packet.ImActEn, packet.ExActEn, packet.ImRctEn, packet.ExRctEn, packet.TlActEn, packet.TlRctEn)
                        if err != nil {
                                glog.Error(err)
                                return err
                        }
                        rows, err := result.RowsAffected()
                        if rows != 1 {
                                glog.Error("expected to affect 1 row, affected %d", rows)
                                return err
                        }
                }
        }
        return nil
}

func GetDummyPackets(n int) Packets {
        pcks := make(Packets, 0)
        start := time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC)
        ip := start
        for i:=0; i<n; i++ {
                ip = ip.Add(time.Minute*5)
                packet := &Packet{Id:int64(n), Timestamp: ip, Voltage: 240.0+rand.Float64(), Frequency: 49.0+rand.Float64(), Status: false, ActiPwr:100*rand.Float64(), AppaPwr:50*rand.Float64(), ReacPwr: 50.0*rand.Float64(), PwrFctr: rand.Float64(), ImActEn: 2.0+rand.Float64(), ExActEn:2.0*rand.Float64(), ImRctEn: 2.0*rand.Float64(), ExRctEn:2.0*rand.Float64(), TlActEn: 2.0*rand.Float64(), TlRctEn: 2.0*rand.Float64()}
                pcks = append(pcks, packet)
        }
        return pcks
}

type Summary struct {
        Hash int64 `json:"hash"`
        Timestamp time.Time `json:"timestamp"`
        VoltageMax float64 `json:"voltagemax"`
        VoltageMin float64 `json:"voltagemin"`
        VoltageAve float64 `json:"voltageave"`
        VoltageExcpt int32 `json:"voltageexcpt"`
        ActiPwrMax float64 `json:"activepwrmax"`
        ActiPwrMin float64 `json:"activepwrmin"`
        ActiPwrAve float64 `json:"activepwrave"`
        AppaPwrMax float64 `json:"apparentpwrmax"`
        AppaPwrMin float64 `json:"apparentpwrmin"`
        AppaPwrAve float64 `json:"apparentpwrave"`
        ReacPwrMax float64 `json:"reactivepwrmax"`
        ReacPwrMin float64 `json:"reactivepwrmin"`
        ReacPwrAve float64 `json:"reactivepwrave"`
        FrequencyMax float64 `json:"frequencymax"`
        FrequencyMin float64 `json:"frequencymin"`
        FrequencyAve float64 `json:"frequencyave"`
        FrequencyExcpt int32 `json:"frequencyexcpt"`
        ImportActiveEnergy float64 `json:"importactiveenergy"`
        ExportActiveEnergy float64 `json:"exportactiveenergy"`
        ImportReactiveEnergy float64 `json:"importreactiveenergy"`
        ExportReactiveEnergy float64 `json:"exportreactiveenergy"`
        TotalActiveEnergy float64 `json:"totalactiveenergy"`
        TotalReactiveEnergy  float64 `json:"totalreactiveenergy"`
}

type Summaries []*Summary

//Summarise Does not check for zero length slices
func (ss Summaries) Summarise(table string) error /**Summary*/ {
        voltageMax := func(s1, s2 *Summary) bool {
                return s1.VoltageMax < s2.VoltageMax // increasing voltage
        }
        voltageMin := func(s1, s2 *Summary) bool {
                return s1.VoltageMin < s2.VoltageMin // increasing voltage
        }
        freqMax := func(s1, s2 *Summary) bool {
                return s1.FrequencyMax < s2.FrequencyMax // increasing frequency
        }
        freqMin := func(s1, s2 *Summary) bool {
                return s1.FrequencyMin < s2.FrequencyMin // increasing frequency
        }
        pwrMax := func(p1, p2 *Summary) bool {
                return p1.ActiPwrMax < p2.ActiPwrMax // increasing power
        }
        pwrMin := func(p1, p2 *Summary) bool {
                return p1.ActiPwrMin < p2.ActiPwrMin // increasing power
        }
        timestamp := func(s1, s2 *Summary) bool {
                return s1.Timestamp.After(s2.Timestamp) // latest first
        }
        OrderedSummaryBy(voltageMax).Sort(ss)
        maxv := ss[len(ss)-1].VoltageMax
        OrderedSummaryBy(voltageMin).Sort(ss)
        minv := ss[0].VoltageMin
        OrderedSummaryBy(freqMax).Sort(ss)
        maxf := ss[len(ss)-1].FrequencyMax
        OrderedSummaryBy(freqMin).Sort(ss)
        minf := ss[len(ss)-1].FrequencyMin
        OrderedSummaryBy(pwrMax).Sort(ss)
        maxp := ss[len(ss)-1].ActiPwrMax
        OrderedSummaryBy(pwrMin).Sort(ss)
        minp := ss[0].ActiPwrMin
        OrderedSummaryBy(timestamp).Sort(ss)
        summary := &Summary{Hash: ss[0].Hash, Timestamp: ss[0].Timestamp, VoltageMax: maxv, VoltageMin: minv, FrequencyMax: maxf, FrequencyMin: minf, ActiPwrMax: maxp, ActiPwrMin: minp, ActiPwrAve: 0.0, ImportActiveEnergy: ss[0].ImportActiveEnergy , ExportActiveEnergy: ss[0].ExportActiveEnergy, ImportReactiveEnergy: ss[0].ImportReactiveEnergy, ExportReactiveEnergy: ss[0].ExportReactiveEnergy, TotalActiveEnergy: ss[0].TotalActiveEnergy, TotalReactiveEnergy: ss[0].TotalReactiveEnergy}
        _, err := PutSummary(table, summary)
        if err != nil {
                glog.Errorf("%v \n", err)
                return err
        }
        return nil
}

//PutSummary persists a Summary into the tablename provided as arg `table`
func PutSummary(table string, hourly *Summary) (uint64, error) {
        db, err := GetDB()
        if err != nil {
                glog.Error(err)
                return 0, err
        }
        result, err := db.Exec("insert into " + table + " (pub_hash, timestamp, voltage_max, voltage_min, voltage_ave, voltage_exceptions, frequency_max, frequency_min, frequency_ave, frequency_exceptions, activepwr_max, activepwr_min, activepwr_ave, import_active_energy, export_active_energy, import_reactive_energy, export_reactive_energy, total_active_energy, total_reactive_energy) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)", hourly.Hash, hourly.Timestamp, hourly.VoltageMax, hourly.VoltageMin, hourly.VoltageAve, hourly.VoltageExcpt, hourly.FrequencyMax, hourly.FrequencyMin, hourly.FrequencyAve, hourly.FrequencyExcpt, hourly.ActiPwrMax, hourly.ActiPwrMin, hourly.ActiPwrAve, hourly.ImportActiveEnergy, hourly.ExportActiveEnergy, hourly.ImportReactiveEnergy, hourly.ExportReactiveEnergy, hourly.TotalActiveEnergy, hourly.TotalReactiveEnergy)
        if err != nil {
                glog.Error(err)
                return 0 , err
        }
        rows, err := result.RowsAffected()
        if rows != 1 {
                glog.Error("expected to affect 1 row, affected %d", rows)
                return uint64(rows) , err
        }
        return uint64(rows), nil
}

func GetSummaries(from, to time.Time, table string) (Summaries, error) {
        db, err := GetDB()
        if err != nil {
                glog.Error(err)
                return nil, err
        }
        pcks := make([]*Summary, 0)
        var query string
        switch table {
        case "hourly":
                query = gethourlies
        case "daily":
                query = getdailies
        default:
                return pcks, fmt.Errorf("No data for %s \n", table)
        }
        rows, err := db.Query(query, from, to)
        if err != nil {
                glog.Errorf("data.GetHourlies %v \n", err)
                return nil, err
        }
        defer rows.Close()
        for rows.Next() {
                pck := &Summary{}
                if err := rows.Scan(&pck.Hash, &pck.Timestamp, &pck.VoltageMax, &pck.VoltageMin, &pck.VoltageAve, &pck.VoltageExcpt, &pck.FrequencyMax, &pck.FrequencyMin, &pck.FrequencyAve, &pck.FrequencyExcpt, &pck.ActiPwrMax, &pck.ActiPwrMin, &pck.ActiPwrAve, &pck.ImportActiveEnergy, &pck.ExportActiveEnergy, &pck.ImportReactiveEnergy, &pck.ExportReactiveEnergy, &pck.TotalActiveEnergy, &pck.TotalReactiveEnergy); err != nil {
                        glog.Errorf("data.GetHourlies %v \n", err)
                        return pcks, fmt.Errorf("No data for hourlies \n")
                }
                //glog.Infof("data.GetSubs appending \n")
                pcks = append(pcks, pck)
        }
        return pcks, nil
}

func GetSummariesByHash(from, to time.Time, table string, pub_hash int64) (Summaries, error) {
        db, err := GetDB()
        if err != nil {
                glog.Error(err)
                return nil, err
        }
        pcks := make([]*Summary, 0)
        var query string
        switch table {
        case "hourly":
                query = gethourliesbyhash
        case "daily":
                query = getdailiesbyhash
        default:
                return pcks, fmt.Errorf("No data for %s \n", table)
        }
        rows, err := db.Query(query, from, to, pub_hash)
        if err != nil {
                glog.Errorf("data.GetHourlies %v \n", err)
                return nil, err
        }
        defer rows.Close()
        for rows.Next() {
                pck := &Summary{}
                if err := rows.Scan(&pck.Hash, &pck.Timestamp, &pck.VoltageMax, &pck.VoltageMin, &pck.VoltageAve, &pck.VoltageExcpt, &pck.FrequencyMax, &pck.FrequencyMin, &pck.FrequencyAve, &pck.FrequencyExcpt, &pck.ActiPwrMax, &pck.ActiPwrMin, &pck.ActiPwrAve, &pck.ImportActiveEnergy, &pck.ExportActiveEnergy, &pck.ImportReactiveEnergy, &pck.ExportReactiveEnergy, &pck.TotalActiveEnergy, &pck.TotalReactiveEnergy); err != nil {
                        glog.Errorf("data.GetHourlies %v \n", err)
                        return pcks, fmt.Errorf("No data for hourlies \n")
                }
                //glog.Infof("data.GetSubs appending \n")
                pcks = append(pcks, pck)
        }
        return pcks, nil
}

func GetStartingSummaryTime() (time.Time, error) {
        db, err := GetDB()
        if err != nil {
                glog.Error(err)
                return time.Now(), err
        }
        rows, err := db.Query(getlastsummarytime)
        if err != nil {
                glog.Errorf("data.GetLastSummaryTime %v \n", err)
                return time.Now(), err
        }
        defer rows.Close()
        pck := &Summary{}
        for rows.Next() {
                if err := rows.Scan(&pck.Hash, &pck.Timestamp); err != nil {
                        glog.Errorf("data.GetLastSummaryTime %v \n", err)
                        return time.Now(), fmt.Errorf("No data for hourlies \n")
                }
        }
        return pck.Timestamp.Round(time.Hour), nil
}

type lessFunc func(p1, p2 *Packet) bool

type multiPacketSorter struct {
        packets Packets
        less []lessFunc
}

// Sort sorts the argument slice according to the less functions passed to OrderedBy.
func (ms *multiPacketSorter) Sort(packets []*Packet) {
        ms.packets = packets
        sort.Sort(ms)
}

// OrderedBy returns a Sorter that sorts using the less functions, in order.
// Call its Sort method to sort the data.
func OrderedBy(less ...lessFunc) *multiPacketSorter {
        return &multiPacketSorter{
                less: less,
        }
}

// Len is part of sort.Interface.
func (ms *multiPacketSorter) Len() int {
        return len(ms.packets)
}

// Swap is part of sort.Interface.
func (ms *multiPacketSorter) Swap(i, j int) {
        ms.packets[i], ms.packets[j] = ms.packets[j], ms.packets[i]
}

// Less is part of sort.Interface. It is implemented by looping along the
// less functions until it finds a comparison that discriminates between
// the two items (one is less than the other). Note that it can call the
// less functions twice per call. We could change the functions to return
// -1, 0, 1 and reduce the number of calls for greater efficiency: an
// exercise for the reader.
func (ms *multiPacketSorter) Less(i, j int) bool {
        p, q := ms.packets[i], ms.packets[j]
        // Try all but the last comparison.
        var k int
        for k = 0; k < len(ms.less)-1; k++ {
                less := ms.less[k]
                switch {
                case less(p, q):
                        // p < q, so we have a decision.
                        return true
                case less(q, p):
                        // p > q, so we have a decision.
                        return false
                }
                // p == q; try the next comparison.
        }
        // All comparisons to here said "equal", so just return whatever
        // the final comparison reports.
        return ms.less[k](p, q)
}

type lessSummaryFunc func(s1, s2 *Summary) bool

type multiSummarySorter struct {
        summaries Summaries
        less []lessSummaryFunc
}

// Sort sorts the argument slice according to the less functions passed to OrderedSummaryBy.
func (ms *multiSummarySorter) Sort(summaries []*Summary) {
        ms.summaries = summaries
        sort.Sort(ms)
}

func OrderedSummaryBy(less ...lessSummaryFunc) *multiSummarySorter {
        return &multiSummarySorter{
                less: less,
        }
}

func (ms *multiSummarySorter) Len() int {
        return len(ms.summaries)
}

func (ms *multiSummarySorter) Swap(i, j int) {
        ms.summaries[i], ms.summaries[j] = ms.summaries[j], ms.summaries[i]
}

func (ms *multiSummarySorter) Less(i, j int) bool {
        p, q := ms.summaries[i], ms.summaries[j]
        // Try all but the last comparison.
        var k int
        for k = 0; k < len(ms.less)-1; k++ {
                less := ms.less[k]
                switch {
                case less(p, q):
                        // p < q, so we have a decision.
                        return true
                case less(q, p):
                        // p > q, so we have a decision.
                        return false
                }
                // p == q; try the next comparison.
        }
        // All comparisons to here said "equal", so just return whatever
        // the final comparison reports.
        return ms.less[k](p, q)
}
