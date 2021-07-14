package entities

import (
	"errors"
	"strings"
	"time"

	"github.com/containers/podman/v3/libpod/define"
	"github.com/containers/podman/v3/pkg/specgen"
	"github.com/containers/podman/v3/pkg/util"
	"github.com/opencontainers/runtime-spec/specs-go"
)

type PodKillOptions struct {
	All    bool
	Latest bool
	Signal string
}

type PodKillReport struct {
	Errs []error
	Id   string //nolint
}

type ListPodsReport struct {
	Cgroup     string
	Containers []*ListPodContainer
	Created    time.Time
	Id         string //nolint
	InfraId    string //nolint
	Name       string
	Namespace  string
	// Network names connected to infra container
	Networks []string
	Status   string
	Labels   map[string]string
}

type ListPodContainer struct {
	Id     string //nolint
	Names  string
	Status string
}

type PodPauseOptions struct {
	All    bool
	Latest bool
}

type PodPauseReport struct {
	Errs []error
	Id   string //nolint
}

type PodunpauseOptions struct {
	All    bool
	Latest bool
}

type PodUnpauseReport struct {
	Errs []error
	Id   string //nolint
}

type PodStopOptions struct {
	All     bool
	Ignore  bool
	Latest  bool
	Timeout int
}

type PodStopReport struct {
	Errs []error
	Id   string //nolint
}

type PodRestartOptions struct {
	All    bool
	Latest bool
}

type PodRestartReport struct {
	Errs []error
	Id   string //nolint
}

type PodStartOptions struct {
	All    bool
	Latest bool
}

type PodStartReport struct {
	Errs []error
	Id   string //nolint
}

type PodRmOptions struct {
	All    bool
	Force  bool
	Ignore bool
	Latest bool
}

type PodRmReport struct {
	Err error
	Id  string //nolint
}

// PodCreateOptions provides all possible options for creating a pod and its infra container
// swagger:model PodCreateOptions
type PodCreateOptions struct {
	CGroupParent       string            `json:"cgroup-parent,omitempty"`
	CreateCommand      []string          `json:"create-command,omitempty"`
	Hostname           string            `json:"hostname,omitempty"`
	Infra              bool              `json:"infra"`
	InfraImage         string            `json:"infra-image,omitempty"`
	InfraName          string            `json:"infra-name,omitempty"`
	InfraCommand       string            `json:"infra-command,omitempty"`
	InfraConmonPidFile string            `json:"infra-conmon-pidfile,omitempty"`
	Labels             map[string]string `json:"labels,omitempty"`
	Name               string            `json:"name,omitempty"`
	NetFlags           *NetFlags
	Net                *NetOptions
	Share              []string `json:"share,omitempty"`
	Pid                string   `json:"pid,omitempty"`
	Cpus               float64  `json:"cpus,omitempty"`
	CpusetCpus         string   `json:"cpuset-Cpus,omitempty"`
	Userns             specgen.Namespace
}

type ContainerCLIOpts struct {
	Annotation        []string `json:"annotations,omitempty"`
	Attach            []string `json:"attach,omitempty"`
	Authfile          string   `json:"authfile,omitempty"`
	BlkIOWeight       string   `json:"blk-io-wt,omitempty"`
	BlkIOWeightDevice []string `json:"cblk-wt-device,omitempty"`
	CapAdd            []string `json:"cap-add,omitempty"`
	CapDrop           []string `json:"cap-drop,omitempty"`
	CgroupNS          string   `json:"cgroup-ns,omitempty"`
	CGroupsMode       string   `json:"cgroup-mode,omitempty"`
	CGroupParent      string   `json:"cgroup-parent,omitempty"`
	CIDFile           string   `json:"cid-file,omitempty"`
	ConmonPIDFile     string   `json:"infra-conmon-pidfile,omitempty"`
	CPUPeriod         uint64   `json:"cpu-period,omitempty"`
	CPUQuota          int64    `json:"cpu-quota,omitempty"`
	CPURTPeriod       uint64   `json:"cpu-rt-period,omitempty"`
	CPURTRuntime      int64    `json:"cpu-rt-runtime,omitempty"`
	CPUShares         uint64   `json:"cpu-shares,omitempty"`
	CPUS              float64  `json:"cpus,omitempty"`
	CPUSetCPUs        string   `json:"cpuset-cpus,omitempty"`
	CPUSetMems        string   `json:"cpu-set-mems,omitempty"`
	Devices           []string `json:"devices,omitempty"`
	DeviceCGroupRule  []string `json:"device-cgroup-rule,omitempty"`
	DeviceReadBPs     []string `json:"device-read-bp,omitempty"`
	DeviceReadIOPs    []string `json:"device-read-io,omitempty"`
	DeviceWriteBPs    []string `json:"device-write-bp,omitempty"`
	DeviceWriteIOPs   []string `json:"device-write-io,omitempty"`
	Entrypoint        *string  `json:"centrypoint,omitempty"`
	Env               []string `json:"env,omitempty"`
	EnvHost           bool     `json:"env-host,omitempty"`
	EnvFile           []string `json:"env-file,omitempty"`
	Expose            []string `json:"expose,omitempty"`
	GIDMap            []string `json:"gid-map,omitempty"`
	GroupAdd          []string `json:"group-add,omitempty"`
	HealthCmd         string   `json:"health-cmd,omitempty"`
	HealthInterval    string   `json:"health-interval,omitempty"`
	HealthRetries     uint     `json:"health-retries,omitempty"`
	HealthStartPeriod string   `json:"health-start-period,omitempty"`
	HealthTimeout     string   `json:"health-timeout,omitempty"`
	Hostname          string   `json:"hostname,omitempty"`
	HTTPProxy         bool
	ImageVolume       string
	Init              bool
	InitContainerType string
	InitPath          string
	Interactive       bool
	IPC               string
	KernelMemory      string
	Label             []string
	LabelFile         []string
	LogDriver         string
	LogOptions        []string
	Memory            string
	MemoryReservation string
	MemorySwap        string
	MemorySwappiness  int64
	Name              string `json:"infra-name,omitempty"`
	NoHealthCheck     bool
	OOMKillDisable    bool
	OOMScoreAdj       int
	Arch              string
	OS                string
	Variant           string
	PID               string
	PIDsLimit         *int64
	Platform          string
	Pod               string `json:"pod,omitempty"`
	PodIDFile         string
	Personality       string
	PreserveFDs       uint
	Privileged        bool
	PublishAll        bool
	Pull              string `json:"pull,omitempty"`
	Quiet             bool
	ReadOnly          bool
	ReadOnlyTmpFS     bool
	Restart           string
	Replace           bool
	Requires          []string
	Rm                bool
	RootFS            bool
	Secrets           []string
	SecurityOpt       []string
	SdNotifyMode      string
	ShmSize           string
	SignaturePolicy   string
	StopSignal        string
	StopTimeout       uint
	StorageOpt        []string
	SubUIDName        string
	SubGIDName        string
	Sysctl            []string
	Systemd           string
	Timeout           uint
	TmpFS             []string
	TTY               bool
	Timezone          string
	Umask             string
	UIDMap            []string
	Ulimit            []string
	User              string
	UserNS            string
	UTS               string
	Mount             []string
	Volume            []string
	VolumesFrom       []string
	Workdir           string
	SeccompPolicy     string
	PidFile           string
	IsInfra           bool `json:"infra,omitempty"`

	Net *NetOptions

	CgroupConf []string
}

type PodCreateReport struct {
	Id string //nolint
}

func (p *PodCreateOptions) CPULimits() *specs.LinuxCPU {
	cpu := &specs.LinuxCPU{}
	hasLimits := false

	if p.Cpus != 0 {
		period, quota := util.CoresToPeriodAndQuota(p.Cpus)
		cpu.Period = &period
		cpu.Quota = &quota
		hasLimits = true
	}
	if p.CpusetCpus != "" {
		cpu.Cpus = p.CpusetCpus
		hasLimits = true
	}
	if !hasLimits {
		return cpu
	}
	return cpu
}

func ToPodSpecGen(s specgen.PodSpecGenerator, p *PodCreateOptions) (*specgen.PodSpecGenerator, error) {
	// Basic Config
	s.Name = p.Name
	s.Hostname = p.Hostname
	s.Labels = p.Labels
	s.NoInfra = !p.Infra
	if len(p.InfraCommand) > 0 {
		s.InfraCommand = strings.Split(p.InfraCommand, " ")
	}
	if len(p.InfraConmonPidFile) > 0 {
		s.InfraConmonPidFile = p.InfraConmonPidFile
	}
	s.InfraImage = p.InfraImage
	s.SharedNamespaces = p.Share
	s.PodCreateCommand = p.CreateCommand

	// Networking config

	if p.Net != nil {
		s.NetNS = p.Net.Network
		s.StaticIP = p.Net.StaticIP
		s.StaticMAC = p.Net.StaticMAC
		s.PortMappings = p.Net.PublishPorts
		s.CNINetworks = p.Net.CNINetworks
		s.NetworkOptions = p.Net.NetworkOptions
		if p.Net.UseImageResolvConf {
			s.NoManageResolvConf = true
		}
		s.DNSServer = p.Net.DNSServers
		s.DNSSearch = p.Net.DNSSearch
		s.DNSOption = p.Net.DNSOptions
		s.NoManageHosts = p.Net.NoHosts
		s.HostAdd = p.Net.AddHosts
	}

	// Cgroup
	s.CgroupParent = p.CGroupParent

	// Resource config
	cpuDat := p.CPULimits()
	if s.ResourceLimits == nil {
		s.ResourceLimits = &specs.LinuxResources{}
		s.ResourceLimits.CPU = &specs.LinuxCPU{}
	}
	if cpuDat != nil {
		s.ResourceLimits.CPU = cpuDat
		if p.Cpus != 0 {
			s.CPUPeriod = *cpuDat.Period
			s.CPUQuota = *cpuDat.Quota
		}
	}
<<<<<<< HEAD
	s.Userns = p.Userns
	return nil
=======
	return &s, nil
>>>>>>> 7a74a7544 (InfraContainer Rework)
}

type PodPruneOptions struct {
	Force bool `json:"force" schema:"force"`
}

type PodPruneReport struct {
	Err error
	Id  string //nolint
}

type PodTopOptions struct {
	// CLI flags.
	ListDescriptors bool
	Latest          bool

	// Options for the API.
	Descriptors []string
	NameOrID    string
}

type PodPSOptions struct {
	CtrNames  bool
	CtrIds    bool
	CtrStatus bool
	Filters   map[string][]string
	Format    string
	Latest    bool
	Namespace bool
	Quiet     bool
	Sort      string
}

type PodInspectOptions struct {
	Latest bool

	// Options for the API.
	NameOrID string

	Format string
}

type PodInspectReport struct {
	*define.InspectPodData
}

// PodStatsOptions are options for the pod stats command.
type PodStatsOptions struct {
	// All - provide stats for all running pods.
	All bool
	// Latest - provide stats for the latest pod.
	Latest bool
}

// PodStatsReport includes pod-resource statistics data.
type PodStatsReport struct {
	CPU           string
	MemUsage      string
	MemUsageBytes string
	Mem           string
	NetIO         string
	BlockIO       string
	PIDS          string
	Pod           string
	CID           string
	Name          string
}

// ValidatePodStatsOptions validates the specified slice and options. Allows
// for sharing code in the front- and the back-end.
func ValidatePodStatsOptions(args []string, options *PodStatsOptions) error {
	num := 0
	if len(args) > 0 {
		num++
	}
	if options.All {
		num++
	}
	if options.Latest {
		num++
	}
	switch num {
	case 0:
		// Podman v1 compat: if nothing's specified get all running
		// pods.
		options.All = true
		return nil
	case 1:
		return nil
	default:
		return errors.New("--all, --latest and arguments cannot be used together")
	}
}
