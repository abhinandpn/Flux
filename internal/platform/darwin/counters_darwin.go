package darwin

/*
#include <string.h>
#include <sys/types.h>
#include <sys/socket.h>
#include <ifaddrs.h>
#include <net/if.h>
#include <net/if_dl.h>
#include <stdlib.h>

struct flux_counter {
	char name[IFNAMSIZ];
	unsigned long long received;
	unsigned long long sent;
	unsigned int flags;
};

static int flux_read_counters(
	struct flux_counter *out,
	int capacity
) {
	struct ifaddrs *addresses = NULL;

	if (getifaddrs(&addresses) != 0) {
		return -1;
	}

	int count = 0;

	for (struct ifaddrs *ifa = addresses;
	     ifa != NULL && count < capacity;
	     ifa = ifa->ifa_next) {

		if (ifa->ifa_addr == NULL ||
		    ifa->ifa_addr->sa_family != AF_LINK ||
		    ifa->ifa_data == NULL) {
			continue;
		}

		struct if_data *data = (struct if_data *)ifa->ifa_data;

		int duplicate = 0;
		for (int i = 0; i < count; i++) {
			if (strcmp(out[i].name, ifa->ifa_name) == 0) {
				duplicate = 1;
				break;
			}
		}
		if (duplicate) {
			continue;
		}

		strlcpy(out[count].name, ifa->ifa_name, IFNAMSIZ);
		out[count].received = data->ifi_ibytes;
		out[count].sent = data->ifi_obytes;
		out[count].flags = ifa->ifa_flags;
		count++;
	}

	freeifaddrs(addresses);
	return count;
}
*/
import "C"

import (
	"context"
	"fmt"
	"time"

	"github.com/abhinandpn/flux/internal/core"
)

// ReadCounters returns cumulative byte counters for macOS interfaces.
func ReadCounters(ctx context.Context) ([]core.CounterSample, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	const capacity = 256
	buffer := make([]C.struct_flux_counter, capacity)

	count := C.flux_read_counters(&buffer[0], C.int(capacity))
	if count < 0 {
		return nil, fmt.Errorf("read macOS network counters failed")
	}

	now := time.Now()
	result := make([]core.CounterSample, 0, int(count))

	for i := 0; i < int(count); i++ {
		item := buffer[i]
		name := C.GoString(&item.name[0])
		if name == "" {
			continue
		}

		result = append(result, core.CounterSample{
			Interface:     name,
			BytesReceived: uint64(item.received),
			BytesSent:     uint64(item.sent),
			Timestamp:     now,
		})
	}

	return result, nil
}
