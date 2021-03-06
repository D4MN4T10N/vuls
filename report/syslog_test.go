package report

import (
	"sort"
	"testing"

	"github.com/future-architect/vuls/models"
)

func TestSyslogWriterEncodeSyslog(t *testing.T) {
	var tests = []struct {
		result           models.ScanResult
		expectedMessages []string
	}{
		{
			result: models.ScanResult{
				ServerName: "teste01",
				Family:     "ubuntu",
				Release:    "16.04",
				IPv4Addrs:  []string{"192.168.0.1", "10.0.2.15"},
				ScannedCves: models.VulnInfos{
					"CVE-2017-0001": models.VulnInfo{
						AffectedPackages: models.PackageStatuses{
							models.PackageStatus{Name: "pkg1"},
							models.PackageStatus{Name: "pkg2"},
						},
					},
					"CVE-2017-0002": models.VulnInfo{
						AffectedPackages: models.PackageStatuses{
							models.PackageStatus{Name: "pkg3"},
							models.PackageStatus{Name: "pkg4"},
						},
						CveContents: models.CveContents{
							models.NVD: models.CveContent{
								Cvss2Score:  5.0,
								Cvss2Vector: "AV:L/AC:L/Au:N/C:N/I:N/A:C",
								CweID:       "CWE-20",
							},
						},
					},
				},
			},
			expectedMessages: []string{
				`server_name="teste01" os_family="ubuntu" os_release="16.04" ipv4_addr="192.168.0.1,10.0.2.15" ipv6_addr="" packages="pkg1,pkg2" cve_id="CVE-2017-0001"`,
				`server_name="teste01" os_family="ubuntu" os_release="16.04" ipv4_addr="192.168.0.1,10.0.2.15" ipv6_addr="" packages="pkg3,pkg4" cve_id="CVE-2017-0002" severity="MEDIUM" cvss_score_v2="5.00" cvss_vector_v2="AV:L/AC:L/Au:N/C:N/I:N/A:C" cwe_id="CWE-20"`,
			},
		},
		{
			result: models.ScanResult{
				ServerName: "teste02",
				Family:     "centos",
				Release:    "6",
				IPv6Addrs:  []string{"2001:0DB8::1"},
				ScannedCves: models.VulnInfos{
					"CVE-2017-0003": models.VulnInfo{
						AffectedPackages: models.PackageStatuses{
							models.PackageStatus{Name: "pkg5"},
						},
						CveContents: models.CveContents{
							models.RedHat: models.CveContent{
								Cvss3Score:  5.0,
								Cvss3Vector: "AV:L/AC:L/Au:N/C:N/I:N/A:C",
								CweID:       "CWE-284",
							},
						},
					},
				},
			},
			expectedMessages: []string{
				`server_name="teste02" os_family="centos" os_release="6" ipv4_addr="" ipv6_addr="2001:0DB8::1" packages="pkg5" cve_id="CVE-2017-0003"`,
			},
		},
	}

	for i, tt := range tests {
		messages := SyslogWriter{}.encodeSyslog(tt.result)
		if len(messages) != len(tt.expectedMessages) {
			t.Fatalf("test: %d, Message Length: expected %d, actual: %d",
				i, len(tt.expectedMessages), len(messages))
		}

		sort.Slice(messages, func(i, j int) bool {
			return messages[i] < messages[j]
		})

		for j, m := range messages {
			e := tt.expectedMessages[j]
			if e != m {
				t.Errorf("test: %d, Messsage %d: expected %s, actual %s", i, j, e, m)
			}
		}
	}
}
