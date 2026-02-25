package main

import (
	"os/signal"
	"syscall"
	"context"
	"os/signal"
	"syscall"
	"context"
	"crypto/tls"
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
)

type CertInfo struct {
	Domain       string
	Issuer       string
	NotBefore    time.Time
	NotAfter     time.Time
	DaysRemaining int
	Status       string
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println(color.CyanString("certwatch - SSL Certificate Watcher"))
		fmt.Println()
		fmt.Println("Usage: certwatch <domain1> <domain2> ...")
		os.Exit(1)
	}

	fmt.Println(color.CyanString("\n=== CERTIFICATE WATCHER ===\n"))

	for _, arg := range os.Args[1:] {
		info, err := checkCertificate(arg)
		if err != nil {
			color.Red("Error checking %s: %v", arg, err)
			continue
		}
		displayCertInfo(info)
	}

	fmt.Println(color.YellowString("\n=== RENEWAL SCHEDULE ==="))
	displayRenewalSchedule()
}

func checkCertificate(domain string) (*CertInfo, error) {
	conn, err := tls.Dial("tcp", domain+":443", nil)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	certs := conn.ConnectionState().PeerCertificates
	if len(certs) == 0 {
		return nil, fmt.Errorf("no certificates found")
	}

	cert := certs[0]
	now := time.Now()
	daysRemaining := int(cert.NotAfter.Sub(now).Hours() / 24)

	status := "OK"
	if daysRemaining < 30 {
		status = "EXPIRING_SOON"
	} else if daysRemaining < 7 {
		status = "CRITICAL"
	} else if daysRemaining < 0 {
		status = "EXPIRED"
	}

	return &CertInfo{
		Domain:        domain,
		Issuer:        cert.Issuer.CommonName,
		NotBefore:     cert.NotBefore,
		NotAfter:      cert.NotAfter,
		DaysRemaining: daysRemaining,
		Status:        status,
	}, nil
}

func displayCertInfo(info *CertInfo) {
	statusColor := color.GreenString
	if info.Status == "EXPIRING_SOON" {
		statusColor = color.YellowString
	} else if info.Status == "CRITICAL" || info.Status == "EXPIRED" {
		statusColor = color.RedString
	}

	fmt.Printf("Domain: %s\n", color.HiWhiteString(info.Domain))
	fmt.Printf("Issuer: %s\n", info.Issuer)
	fmt.Printf("Valid From: %s\n", info.NotBefore.Format("2006-01-02"))
	fmt.Printf("Valid Until: %s\n", info.NotAfter.Format("2006-01-02"))
	fmt.Printf("Days Remaining: %s\n", statusColor(fmt.Sprintf("%d days", info.DaysRemaining)))
	fmt.Printf("Status: %s\n\n", statusColor(info.Status))
}

func displayRenewalSchedule() {
	fmt.Println(color.CyanString("Recommended cron jobs for auto-renewal:"))
	fmt.Println()
	fmt.Println("30 days before expiry:")
	fmt.Println("  0 0 1 * * certbot renew --force-renewal")
	fmt.Println()
	fmt.Println("7 days before expiry:")
	fmt.Println("  0 0 1 * * certbot renew")
	fmt.Println()
	fmt.Println("Daily check:")
	fmt.Println("  0 */6 * * * certwatch <domains> | mail -s 'Cert Alert' admin@example.com")
}