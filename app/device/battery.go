package device

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const powerSupplyPath = "/sys/devices/platform/pmic_battery.1/power_supply/mc13892_bat"

func GetBatteryLevel() int {
	capacityFile := fmt.Sprintf("%s/capacity", powerSupplyPath)

	percentage, err := os.ReadFile(capacityFile)
	if err != nil {
		return 0
	}

	percentageNumber, err := strconv.Atoi(strings.TrimSuffix(string(percentage), "\n"))
	if err != nil {
		return 0
	}

	return percentageNumber
}
