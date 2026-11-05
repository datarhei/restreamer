package restream

import (
	"fmt"
	"testing"
	"time"

	"github.com/datarhei/core/v16/ffmpeg"
	"github.com/datarhei/core/v16/internal/testhelper"
	"github.com/datarhei/core/v16/net"
	"github.com/datarhei/core/v16/restream/app"
	"github.com/datarhei/core/v16/restream/replace"

	"github.com/stretchr/testify/require"
)

func getDummyRestreamer(portrange net.Portranger, validatorIn, validatorOut ffmpeg.Validator, replacer replace.Replacer) (Restreamer, error) {
	binary, err := testhelper.BuildBinary("ffmpeg", "../internal/testhelper")
	if err != nil {
		return nil, fmt.Errorf("failed to build helper program: %w", err)
	}

	ffmpeg, err := ffmpeg.New(ffmpeg.Config{
		Binary:           binary,
		LogHistoryLength: 3,
		Portrange:        portrange,
		ValidatorInput:   validatorIn,
		ValidatorOutput:  validatorOut,
	})
	if err != nil {
		return nil, err
	}

	rs, err := New(Config{
		FFmpeg:  ffmpeg,
		Replace: replacer,
	})
	if err != nil {
		return nil, err
	}

	return rs, nil
}

func getDummyProcess() *app.Config {
	return &app.Config{
		ID: "process",
		Input: []app.ConfigIO{
			{
				ID:      "in",
				Address: "testsrc=size=1280x720:rate=25",
				Options: []string{
					"-f",
					"lavfi",
					"-re",
				},
			},
		},
		Output: []app.ConfigIO{
			{
				ID:      "out",
				Address: "-",
				Options: []string{
					"-codec",
					"copy",
					"-f",
					"null",
				},
			},
		},
		Options: []string{
			"-loglevel",
			"info",
		},
		Reconnect:      true,
		ReconnectDelay: 10,
		Autostart:      false,
		StaleTimeout:   0,
	}
}

func TestAddProcess(t *testing.T) {
	rs, err := getDummyRestreamer(nil, nil, nil, nil)
	require.NoError(t, err)

	process := getDummyProcess()
	require.NotNil(t, process)

	_, err = rs.GetProcess(process.ID)
	require.NotEqual(t, nil, err, "Unset process found (%s)", process.ID)

	err = rs.AddProcess(process)
	require.Equal(t, nil, err, "Failed to add process (%s)", err)

	_, err = rs.GetProcess(process.ID)
	require.Equal(t, nil, err, "Set process not found (%s)", process.ID)

	state, _ := rs.GetProcessState(process.ID)
	require.Equal(t, "stop", state.Order, "Process should be stopped")
}

func TestAutostartProcess(t *testing.T) {
	rs, err := getDummyRestreamer(nil, nil, nil, nil)
	require.NoError(t, err)

	process := getDummyProcess()
	process.Autostart = true

	rs.AddProcess(process)

	state, _ := rs.GetProcessState(process.ID)
	require.Equal(t, "start", state.Order, "Process should be started")

	rs.StopProcess(process.ID)
}

func TestAddInvalidProcess(t *testing.T) {
	rs, err := getDummyRestreamer(nil, nil, nil, nil)
	require.NoError(t, err)

	// Invalid process ID
	process := getDummyProcess()
	process.ID = ""

	err = rs.AddProcess(process)
	require.NotEqual(t, nil, err, "Succeeded to add process without ID")

	// Invalid input ID
	process = getDummyProcess()
	process.Input[0].ID = ""

	err = rs.AddProcess(process)
	require.NotEqual(t, nil, err, "Succeeded to add process input without ID")

	// Invalid input address
	process = getDummyProcess()
	process.Input[0].Address = ""

	err = rs.AddProcess(process)
	require.NotEqual(t, nil, err, "Succeeded to add process input without address")

	// Duplicate input ID
	process = getDummyProcess()
	process.Input = append(process.Input, process.Input[0])

	err = rs.AddProcess(process)
	require.NotEqual(t, nil, err, "Succeeded to add process input with duplicate ID")

	// No inputs
	process = getDummyProcess()
	process.Input = nil

	err = rs.AddProcess(process)
	require.NotEqual(t, nil, err, "Succeeded to add process without inputs")

	// Invalid output ID
	process = getDummyProcess()
	process.Output[0].ID = ""

	err = rs.AddProcess(process)
	require.NotEqual(t, nil, err, "Succeeded to add process output without ID")

	// Invalid output address
	process = getDummyProcess()
	process.Output[0].Address = ""

	err = rs.AddProcess(process)
	require.NotEqual(t, nil, err, "Succeeded to add process output without address")

	// Duplicate output ID
	process = getDummyProcess()
	process.Output = append(process.Output, process.Output[0])

	err = rs.AddProcess(process)
	require.NotEqual(t, nil, err, "Succeeded to add process output with duplicate ID")

	// No outputs
	process = getDummyProcess()
	process.Output = nil

	err = rs.AddProcess(process)
	require.NotEqual(t, nil, err, "Succeeded to add process without outputs")
}

func TestRemoveProcess(t *testing.T) {
	rs, err := getDummyRestreamer(nil, nil, nil, nil)
	require.NoError(t, err)

	process := getDummyProcess()

	err = rs.AddProcess(process)
	require.Equal(t, nil, err, "Failed to add process (%s)", err)

	err = rs.DeleteProcess(process.ID)
	require.Equal(t, nil, err, "Set process not found (%s)", process.ID)

	_, err = rs.GetProcess(process.ID)
	require.NotEqual(t, nil, err, "Unset process found (%s)", process.ID)
}

func TestUpdateProcess(t *testing.T) {
	rs, err := getDummyRestreamer(nil, nil, nil, nil)
	require.NoError(t, err)

	process1 := getDummyProcess()
	require.NotNil(t, process1)
	process1.ID = "process1"

	process2 := getDummyProcess()
	require.NotNil(t, process2)
	process2.ID = "process2"

	err = rs.AddProcess(process1)
	require.Equal(t, nil, err)

	err = rs.AddProcess(process2)
	require.Equal(t, nil, err)

	process, err := rs.GetProcess(process2.ID)
	require.NoError(t, err)

	createdAt := process.CreatedAt
	updatedAt := process.UpdatedAt

	time.Sleep(2 * time.Second)

	process3 := getDummyProcess()
	require.NotNil(t, process3)
	process3.ID = "process2"

	err = rs.UpdateProcess("process1", process3)
	require.Error(t, err)

	process3.ID = "process3"
	err = rs.UpdateProcess("process1", process3)
	require.NoError(t, err)

	_, err = rs.GetProcess(process1.ID)
	require.Error(t, err)

	process, err = rs.GetProcess(process3.ID)
	require.NoError(t, err)

	require.NotEqual(t, createdAt, process.CreatedAt) // this should be equal, but will require a major version jump
	require.NotEqual(t, updatedAt, process.UpdatedAt)
}

func TestGetProcess(t *testing.T) {
	rs, err := getDummyRestreamer(nil, nil, nil, nil)
	require.NoError(t, err)

	process1 := getDummyProcess()
	process1.ID = "foo_aaa_1"
	process1.Reference = "foo_aaa_1"
	process2 := getDummyProcess()
	process2.ID = "bar_bbb_2"
	process2.Reference = "bar_bbb_2"
	process3 := getDummyProcess()
	process3.ID = "foo_ccc_3"
	process3.Reference = "foo_ccc_3"
	process4 := getDummyProcess()
	process4.ID = "bar_ddd_4"
	process4.Reference = "bar_ddd_4"

	rs.AddProcess(process1)
	rs.AddProcess(process2)
	rs.AddProcess(process3)
	rs.AddProcess(process4)

	_, err = rs.GetProcess(process1.ID)
	require.Equal(t, nil, err)

	list := rs.GetProcessIDs("", "")
	require.Len(t, list, 4)
	require.ElementsMatch(t, []string{"foo_aaa_1", "bar_bbb_2", "foo_ccc_3", "bar_ddd_4"}, list)

	list = rs.GetProcessIDs("foo_*", "")
	require.Len(t, list, 2)
	require.ElementsMatch(t, []string{"foo_aaa_1", "foo_ccc_3"}, list)

	list = rs.GetProcessIDs("bar_*", "")
	require.Len(t, list, 2)
	require.ElementsMatch(t, []string{"bar_bbb_2", "bar_ddd_4"}, list)

	list = rs.GetProcessIDs("*_bbb_*", "")
	require.Len(t, list, 1)
	require.ElementsMatch(t, []string{"bar_bbb_2"}, list)

	list = rs.GetProcessIDs("", "foo_*")
	require.Len(t, list, 2)
	require.ElementsMatch(t, []string{"foo_aaa_1", "foo_ccc_3"}, list)

	list = rs.GetProcessIDs("", "bar_*")
	require.Len(t, list, 2)
	require.ElementsMatch(t, []string{"bar_bbb_2", "bar_ddd_4"}, list)

	list = rs.GetProcessIDs("", "*_bbb_*")
	require.Len(t, list, 1)
	require.ElementsMatch(t, []string{"bar_bbb_2"}, list)
}

func TestStartProcess(t *testing.T) {
	rs, err := getDummyRestreamer(nil, nil, nil, nil)
	require.NoError(t, err)

	process := getDummyProcess()

	rs.AddProcess(process)

	err = rs.StartProcess("foobar")
	require.NotEqual(t, nil, err, "shouldn't be able to start non-existing process")

	err = rs.StartProcess(process.ID)
	require.Equal(t, nil, err, "should be able to start existing process")

	state, _ := rs.GetProcessState(process.ID)
	require.Equal(t, "start", state.Order, "Process should be started")

	err = rs.StartProcess(process.ID)
	require.Equal(t, nil, err, "should be able to start already running process")

	state, _ = rs.GetProcessState(process.ID)
	require.Equal(t, "start", state.Order, "Process should be started")

	rs.StopProcess(process.ID)
}

func TestStopProcess(t *testing.T) {
	rs, err := getDummyRestreamer(nil, nil, nil, nil)
	require.NoError(t, err)

	process := getDummyProcess()

	rs.AddProcess(process)
	rs.StartProcess(process.ID)

	err = rs.StopProcess("foobar")
	require.NotEqual(t, nil, err, "shouldn't be able to stop non-existing process")

	err = rs.StopProcess(process.ID)
	require.Equal(t, nil, err, "should be able to stop existing running process")

	state, _ := rs.GetProcessState(process.ID)
	require.Equal(t, "stop", state.Order, "Process should be stopped")

	err = rs.StopProcess(process.ID)
	require.Equal(t, nil, err, "should be able to stop already stopped process")

	state, _ = rs.GetProcessState(process.ID)
	require.Equal(t, "stop", state.Order, "Process should be stopped")
}

func TestRestartProcess(t *testing.T) {
	rs, err := getDummyRestreamer(nil, nil, nil, nil)
	require.NoError(t, err)

	process := getDummyProcess()

	rs.AddProcess(process)

	err = rs.RestartProcess("foobar")
	require.NotEqual(t, nil, err, "shouldn't be able to restart non-existing process")

	err = rs.RestartProcess(process.ID)
	require.Equal(t, nil, err, "should be able to restart existing stopped process")

	state, _ := rs.GetProcessState(process.ID)
	require.Equal(t, "stop", state.Order, "Process should be stopped")

	rs.StartProcess(process.ID)

	state, _ = rs.GetProcessState(process.ID)
	require.Equal(t, "start", state.Order, "Process should be started")

	rs.StopProcess(process.ID)
}

func TestReloadProcess(t *testing.T) {
	rs, err := getDummyRestreamer(nil, nil, nil, nil)
	require.NoError(t, err)

	process := getDummyProcess()

	rs.AddProcess(process)

	err = rs.ReloadProcess("foobar")
	require.NotEqual(t, nil, err, "shouldn't be able to reload non-existing process")

	err = rs.ReloadProcess(process.ID)
	require.Equal(t, nil, err, "should be able to reload existing stopped process")

	state, _ := rs.GetProcessState(process.ID)
	require.Equal(t, "stop", state.Order, "Process should be stopped")

	rs.StartProcess(process.ID)

	state, _ = rs.GetProcessState(process.ID)
	require.Equal(t, "start", state.Order, "Process should be started")

	err = rs.ReloadProcess(process.ID)
	require.Equal(t, nil, err, "should be able to reload existing process")

	state, _ = rs.GetProcessState(process.ID)
	require.Equal(t, "start", state.Order, "Process should be started")

	rs.StopProcess(process.ID)
}

func TestProbeProcess(t *testing.T) {
	rs, err := getDummyRestreamer(nil, nil, nil, nil)
	require.NoError(t, err)

	process := getDummyProcess()

	rs.AddProcess(process)

	probe := rs.ProbeWithTimeout(process.ID, 5*time.Second)

	require.Equal(t, 3, len(probe.Streams))
}

func TestProcessMetadata(t *testing.T) {
	rs, err := getDummyRestreamer(nil, nil, nil, nil)
	require.NoError(t, err)

	process := getDummyProcess()

	rs.AddProcess(process)

	data, _ := rs.GetProcessMetadata(process.ID, "foobar")
	require.Equal(t, nil, data, "nothing should be stored under the key")

	rs.SetProcessMetadata(process.ID, "foobar", process)

	data, _ = rs.GetProcessMetadata(process.ID, "foobar")
	require.NotEqual(t, nil, data, "there should be something stored under the key")

	p := data.(*app.Config)

	require.Equal(t, process.ID, p.ID, "failed to retrieve stored data")
}

func TestLog(t *testing.T) {
	rs, err := getDummyRestreamer(nil, nil, nil, nil)
	require.NoError(t, err)

	process := getDummyProcess()

	rs.AddProcess(process)

	_, err = rs.GetProcessLog("foobar")
	require.Error(t, err)

	log, err := rs.GetProcessLog(process.ID)
	require.NoError(t, err)
	require.Equal(t, 0, len(log.Prelude))
	require.Equal(t, 0, len(log.Log))

	rs.StartProcess(process.ID)

	time.Sleep(3 * time.Second)

	log, _ = rs.GetProcessLog(process.ID)

	require.NotEqual(t, 0, len(log.Prelude))
	require.NotEqual(t, 0, len(log.Log))

	rs.StopProcess(process.ID)

	log, _ = rs.GetProcessLog(process.ID)

	require.NotEqual(t, 0, len(log.Prelude))
	require.NotEqual(t, 0, len(log.Log))
}

func TestLogTransfer(t *testing.T) {
	rs, err := getDummyRestreamer(nil, nil, nil, nil)
	require.NoError(t, err)

	process := getDummyProcess()

	err = rs.AddProcess(process)
	require.NoError(t, err)

	rs.StartProcess(process.ID)
	time.Sleep(3 * time.Second)
	rs.StopProcess(process.ID)

	rs.StartProcess(process.ID)
	rs.StopProcess(process.ID)

	log, _ := rs.GetProcessLog(process.ID)

	require.Equal(t, 1, len(log.History))

	err = rs.UpdateProcess(process.ID, process)
	require.NoError(t, err)

	log, _ = rs.GetProcessLog(process.ID)

	require.Equal(t, 1, len(log.History))
}

func TestPlayoutNoRange(t *testing.T) {
	rs, err := getDummyRestreamer(nil, nil, nil, nil)
	require.NoError(t, err)

	process := getDummyProcess()

	process.Input[0].Address = "playout:" + process.Input[0].Address

	rs.AddProcess(process)

	_, err = rs.GetPlayout("foobar", process.Input[0].ID)
	require.NotEqual(t, nil, err, "playout of non-existing process should error")

	_, err = rs.GetPlayout(process.ID, "foobar")
	require.NotEqual(t, nil, err, "playout of non-existing input should error")

	addr, _ := rs.GetPlayout(process.ID, process.Input[0].ID)
	require.Equal(t, 0, len(addr), "the playout address should be empty if no port range is given")
}

func TestPlayoutRange(t *testing.T) {
	portrange, err := net.NewPortrange(3000, 3001)
	require.NoError(t, err)

	rs, err := getDummyRestreamer(portrange, nil, nil, nil)
	require.NoError(t, err)

	process := getDummyProcess()

	process.Input[0].Address = "playout:" + process.Input[0].Address

	rs.AddProcess(process)

	_, err = rs.GetPlayout("foobar", process.Input[0].ID)
	require.NotEqual(t, nil, err, "playout of non-existing process should error")

	_, err = rs.GetPlayout(process.ID, "foobar")
	require.NotEqual(t, nil, err, "playout of non-existing input should error")

	addr, _ := rs.GetPlayout(process.ID, process.Input[0].ID)
	require.NotEqual(t, 0, len(addr), "the playout address should not be empty if a port range is given")
	require.Equal(t, "127.0.0.1:3000", addr, "the playout address should be 127.0.0.1:3000")
}

func TestAddressReference(t *testing.T) {
	rs, err := getDummyRestreamer(nil, nil, nil, nil)
	require.NoError(t, err)

	process1 := getDummyProcess()
	process2 := getDummyProcess()

	process2.ID = "process2"

	rs.AddProcess(process1)

	process2.Input[0].Address = "#process:foobar=out"

	err = rs.AddProcess(process2)
	require.NotEqual(t, nil, err, "shouldn't resolve invalid reference")

	process2.Input[0].Address = "#process2:output=out"

	err = rs.AddProcess(process2)
	require.NotEqual(t, nil, err, "shouldn't resolve invalid reference")

	process2.Input[0].Address = "#process:output=foobar"

	err = rs.AddProcess(process2)
	require.NotEqual(t, nil, err, "shouldn't resolve invalid reference")

	process2.Input[0].Address = "#process:output=out"

	err = rs.AddProcess(process2)
	require.Equal(t, nil, err, "should resolve reference")
}

func TestConfigValidation(t *testing.T) {
	rsi, err := getDummyRestreamer(nil, nil, nil, nil)
	require.NoError(t, err)

	rs := rsi.(*restream)

	config := getDummyProcess()

	_, err = rs.validateConfig(config)
	require.NoError(t, err)

	config.Input = []app.ConfigIO{}
	_, err = rs.validateConfig(config)
	require.Error(t, err)

	config = getDummyProcess()
	config.Input[0].ID = ""
	_, err = rs.validateConfig(config)
	require.Error(t, err)

	config = getDummyProcess()
	config.Input[0].Address = ""
	_, err = rs.validateConfig(config)
	require.Error(t, err)

	config = getDummyProcess()
	config.Output = []app.ConfigIO{}
	_, err = rs.validateConfig(config)
	require.Error(t, err)

	config = getDummyProcess()
	config.Output[0].ID = ""
	_, err = rs.validateConfig(config)
	require.Error(t, err)

	config = getDummyProcess()
	config.Output[0].Address = ""
	_, err = rs.validateConfig(config)
	require.Error(t, err)
}

func TestConfigValidationFFmpeg(t *testing.T) {
	valIn, err := ffmpeg.NewValidator([]string{"^https?://"}, nil)
	require.NoError(t, err)

	valOut, err := ffmpeg.NewValidator([]string{"^https?://", "^rtmp://"}, nil)
	require.NoError(t, err)

	rsi, err := getDummyRestreamer(nil, valIn, valOut, nil)
	require.NoError(t, err)

	rs := rsi.(*restream)

	config := getDummyProcess()

	_, err = rs.validateConfig(config)
	require.Error(t, err)

	config.Input[0].Address = "http://stream.example.com/master.m3u8"
	config.Output[0].Address = "http://stream.example.com/master2.m3u8"

	_, err = rs.validateConfig(config)
	require.NoError(t, err)

	config.Output[0].Address = "[f=flv]http://stream.example.com/master2.m3u8"
	_, err = rs.validateConfig(config)
	require.NoError(t, err)

	config.Output[0].Address = "[f=hls]http://stream.example.com/master2.m3u8|[f=flv]rtmp://stream.example.com/stream"
	_, err = rs.validateConfig(config)
	require.NoError(t, err)
}

func TestOutputAddressValidation(t *testing.T) {
	rsi, err := getDummyRestreamer(nil, nil, nil, nil)
	require.NoError(t, err)

	rs := rsi.(*restream)

	type res struct {
		path string
		err  bool
	}

	paths := map[string]res{
		"/dev/null":                            {"file:/dev/null", false},
		"/dev/../etc/passwd":                   {"/etc/passwd", true},
		"/dev/fb0":                             {"file:/dev/fb0", false},
		"/etc/passwd":                          {"/etc/passwd", true},
		"/core/data/../../etc/passwd":          {"/etc/passwd", true},
		"/core/data/./etc/passwd":              {"file:/core/data/etc/passwd", false},
		"file:/core/data/foobar":               {"file:/core/data/foobar", false},
		"http://example.com":                   {"http://example.com", false},
		"-":                                    {"pipe:", false},
		"/core/data/foobar|http://example.com": {"file:/core/data/foobar|http://example.com", false},
		"/core/data/foobar|/etc/passwd":        {"/core/data/foobar|/etc/passwd", true},
		"[f=mpegts]udp://10.0.1.255:1234/":     {"[f=mpegts]udp://10.0.1.255:1234/", false},
		"[f=null]-|[f=null]-":                  {"[f=null]pipe:|[f=null]pipe:", false},
		"[onfail=ignore]/core/data/archive-20121107.mkv|[f=mpegts]udp://10.0.1.255:1234/": {"[onfail=ignore]file:/core/data/archive-20121107.mkv|[f=mpegts]udp://10.0.1.255:1234/", false},
	}

	for path, r := range paths {
		path, _, err := rs.validateOutputAddress(path, "/core/data")

		if r.err {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
		}

		require.Equal(t, r.path, path)
	}
}

func TestMetadata(t *testing.T) {
	rs, err := getDummyRestreamer(nil, nil, nil, nil)
	require.NoError(t, err)

	process := getDummyProcess()

	data, _ := rs.GetMetadata("foobar")
	require.Equal(t, nil, data, "nothing should be stored under the key")

	rs.SetMetadata("foobar", process)

	data, _ = rs.GetMetadata("foobar")
	require.NotEqual(t, nil, data, "there should be something stored under the key")

	p := data.(*app.Config)

	require.Equal(t, process.ID, p.ID, "failed to retrieve stored data")
}

func TestReplacer(t *testing.T) {
	replacer := replace.New()

	replacer.RegisterTemplateFunc("diskfs", func(config *app.Config, section string) string {
		return "/mnt/diskfs"
	}, nil)

	replacer.RegisterTemplateFunc("fs:disk", func(config *app.Config, section string) string {
		return "/mnt/diskfs"
	}, nil)

	replacer.RegisterTemplateFunc("memfs", func(config *app.Config, section string) string {
		return "http://localhost/mnt/memfs"
	}, nil)

	replacer.RegisterTemplateFunc("fs:mem", func(config *app.Config, section string) string {
		return "http://localhost/mnt/memfs"
	}, nil)

	replacer.RegisterTemplateFunc("rtmp", func(config *app.Config, section string) string {
		return "rtmp://localhost/app/{name}?token=foobar"
	}, nil)

	replacer.RegisterTemplateFunc("srt", func(config *app.Config, section string) string {
		template := "srt://localhost:6000?mode=caller&transtype=live&latency={latency}&streamid={name}"
		if section == "output" {
			template += ",mode:publish"
		} else {
			template += ",mode:request"
		}
		template += ",token:abcfoobar&passphrase=secret"

		return template
	}, map[string]string{
		"latency": "20000", // 20 milliseconds, FFmpeg requires microseconds
	})

	rsi, err := getDummyRestreamer(nil, nil, nil, replacer)
	require.NoError(t, err)

	process := &app.Config{
		ID:        "314159265359",
		Reference: "refref",
		Input: []app.ConfigIO{
			{
				ID:      "in_{processid}_{reference}",
				Address: "input:{inputid}_process:{processid}_reference:{reference}_diskfs:{diskfs}/disk.txt_memfs:{memfs}/mem.txt_fsdisk:{fs:disk}/fsdisk.txt_fsmem:{fs:mem}/fsmem.txt_rtmp:{rtmp,name=pmtr}_srt:{srt,name=trs}_rtmp:{rtmp,name=$inputid}",
				Options: []string{
					"-f",
					"lavfi",
					"-re",
					"input:{inputid}",
					"process:{processid}",
					"reference:{reference}",
					"diskfs:{diskfs}/disk.txt",
					"memfs:{memfs}/mem.txt",
					"fsdisk:{fs:disk}/fsdisk.txt",
					"fsmem:{fs:mem}/$inputid.txt",
				},
			},
		},
		Output: []app.ConfigIO{
			{
				ID:      "out_{processid}_{reference}",
				Address: "output:{outputid}_process:{processid}_reference:{reference}_diskfs:{diskfs}/disk.txt_memfs:{memfs}/mem.txt_fsdisk:{fs:disk}/fsdisk.txt_fsmem:{fs:mem}/fsmem.txt_rtmp:{rtmp,name=$processid}_srt:{srt,name=$reference,latency=42}_rtmp:{rtmp,name=$outputid}",
				Options: []string{
					"-codec",
					"copy",
					"-f",
					"null",
					"output:{outputid}",
					"process:{processid}",
					"reference:{reference}",
					"diskfs:{diskfs}/disk.txt",
					"memfs:{memfs}/mem.txt",
					"fsdisk:{fs:disk}/fsdisk.txt",
					"fsmem:{fs:mem}/$outputid.txt",
				},
				Cleanup: []app.ConfigIOCleanup{
					{
						Pattern:       "pattern_{outputid}_{processid}_{reference}_{rtmp,name=$outputid}",
						MaxFiles:      0,
						MaxFileAge:    0,
						PurgeOnDelete: false,
					},
				},
			},
		},
		Options: []string{
			"-loglevel",
			"info",
			"{diskfs}/foobar_on_disk.txt",
			"{memfs}/foobar_in_mem.txt",
			"{fs:disk}/foobar_on_disk_aswell.txt",
			"{fs:mem}/foobar_in_mem_aswell.txt",
		},
		Reconnect:      true,
		ReconnectDelay: 10,
		Autostart:      false,
		StaleTimeout:   0,
	}

	err = rsi.AddProcess(process)
	require.NoError(t, err)

	rs := rsi.(*restream)

	process = &app.Config{
		ID:        "314159265359",
		Reference: "refref",
		FFVersion: "^4.0.2",
		Input: []app.ConfigIO{
			{
				ID:      "in_314159265359_refref",
				Address: "input:in_314159265359_refref_process:314159265359_reference:refref_diskfs:/mnt/diskfs/disk.txt_memfs:http://localhost/mnt/memfs/mem.txt_fsdisk:/mnt/diskfs/fsdisk.txt_fsmem:http://localhost/mnt/memfs/fsmem.txt_rtmp:rtmp://localhost/app/pmtr?token=foobar_srt:srt://localhost:6000?mode=caller&transtype=live&latency=20000&streamid=trs,mode:request,token:abcfoobar&passphrase=secret_rtmp:rtmp://localhost/app/in_314159265359_refref?token=foobar",
				Options: []string{
					"-f",
					"lavfi",
					"-re",
					"input:in_314159265359_refref",
					"process:314159265359",
					"reference:refref",
					"diskfs:/mnt/diskfs/disk.txt",
					"memfs:http://localhost/mnt/memfs/mem.txt",
					"fsdisk:/mnt/diskfs/fsdisk.txt",
					"fsmem:http://localhost/mnt/memfs/$inputid.txt",
				},
				Cleanup: []app.ConfigIOCleanup{},
			},
		},
		Output: []app.ConfigIO{
			{
				ID:      "out_314159265359_refref",
				Address: "output:out_314159265359_refref_process:314159265359_reference:refref_diskfs:/mnt/diskfs/disk.txt_memfs:http://localhost/mnt/memfs/mem.txt_fsdisk:/mnt/diskfs/fsdisk.txt_fsmem:http://localhost/mnt/memfs/fsmem.txt_rtmp:rtmp://localhost/app/314159265359?token=foobar_srt:srt://localhost:6000?mode=caller&transtype=live&latency=42&streamid=refref,mode:publish,token:abcfoobar&passphrase=secret_rtmp:rtmp://localhost/app/out_314159265359_refref?token=foobar",
				Options: []string{
					"-codec",
					"copy",
					"-f",
					"null",
					"output:out_314159265359_refref",
					"process:314159265359",
					"reference:refref",
					"diskfs:/mnt/diskfs/disk.txt",
					"memfs:http://localhost/mnt/memfs/mem.txt",
					"fsdisk:/mnt/diskfs/fsdisk.txt",
					"fsmem:http://localhost/mnt/memfs/$outputid.txt",
				},
				Cleanup: []app.ConfigIOCleanup{
					{
						Pattern:       "pattern_out_314159265359_refref_314159265359_refref_{rtmp,name=$outputid}",
						MaxFiles:      0,
						MaxFileAge:    0,
						PurgeOnDelete: false,
					},
				},
			},
		},
		Options: []string{
			"-loglevel",
			"info",
			"/mnt/diskfs/foobar_on_disk.txt",
			"{memfs}/foobar_in_mem.txt",
			"/mnt/diskfs/foobar_on_disk_aswell.txt",
			"http://localhost/mnt/memfs/foobar_in_mem_aswell.txt",
		},
		Reconnect:      true,
		ReconnectDelay: 10,
		Autostart:      false,
		StaleTimeout:   0,
	}

	require.Equal(t, process, rs.tasks["314159265359"].config)
}

func TestProcessLimit(t *testing.T) {
	rsi, err := getDummyRestreamer(nil, nil, nil, nil)
	require.NoError(t, err)

	process := getDummyProcess()
	process.LimitCPU = 61
	process.LimitMemory = 42
	process.Autostart = false

	err = rsi.AddProcess(process)
	require.NoError(t, err)

	rs := rsi.(*restream)

	task, ok := rs.tasks[process.ID]
	require.True(t, ok)

	status := task.ffmpeg.Status()

	require.Equal(t, float64(61), status.CPU.Limit)
	require.Equal(t, uint64(42), status.Memory.Limit)
}
