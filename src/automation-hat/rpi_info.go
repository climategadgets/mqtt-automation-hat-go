package automation_hat

import (
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"strings"
)

var procCpuinfo = "/proc/cpuinfo"
var fieldRevision = "Revision"

// Returns revision if the application is currently running on Raspberry Pi, otherwise nil
// VT: FIXME: We may need a more detailed breakdown - but later
//
// NOTE: "Hardware" entry must not be used, see https://www.raspberrypi.org/documentation/hardware/raspberrypi/revision-codes/README.md
// NOTE: Revision to model breakdown: https://elinux.org/RPi_HardwareHistory
func GetRaspberryPiRevision() *string {

	cpuinfo, err := ioutil.ReadFile(procCpuinfo)

	if err != nil {
		log.Warn(err)
		log.Warn("no " + procCpuinfo + ", likely not running on Pi (nor on UNIX)")
		return nil
	}

	lines := strings.Split(string(cpuinfo), "\n")

	for _, line := range lines {

		if strings.HasPrefix(line, fieldRevision) {
			fields := strings.Fields(line)

			log.Infof("Raspberry Pi: "+fieldRevision+": %v", fields[2])
			return &fields[2]
		}
	}

	log.Warn("no " + fieldRevision + " line in " + procCpuinfo + ", likely not running on Pi")
	return nil
}
