/*
Copyright © 2026 @mdxabu
*/
package cmd

import (
	"os"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"fmt"
)



var rootCmd = &cobra.Command{
	Use:   "twich",
	Short: "Twitch Chat TUI in Go",
	Long: `Twitch Chat TUI in Go`,
	
	Run: func(cmd *cobra.Command, args []string) {
		asciiBanner := []string{
			"████████ ██     ██ ██  ██████ ██   ██ ",
			"   ██    ██     ██ ██ ██      ██   ██ ",
			"   ██    ██  █  ██ ██ ██      ███████ ",
			"   ██    ██ ███ ██ ██ ██      ██   ██ ",
			"   ██     ███ ███  ██  ██████ ██   ██ ",
			"                                      ",
		}

		// Linear gradient from violet (#8B00FF) to red (#FF0000)
		startR, startG, startB := 139, 0, 255   // Violet
		endR, endG, endB := 255, 0, 0           // Red
		lines := len(asciiBanner)

		for i, line := range asciiBanner {
			// Calculate color for this line
			r := startR + (endR-startR)*i/(lines-1)
			g := startG + (endG-startG)*i/(lines-1)
			b := startB + (endB-startB)*i/(lines-1)
			c := color.New(color.FgHiWhite).Add(color.Attribute(0))
			c.Set()
			// Use fatih/color's RGB support if available (since v1.13.0)
			// color.FgRGB is not available; use ANSI escape codes for 24-bit color
			// Format: \x1b[38;2;<r>;<g>;<b>m
			ansiColor := fmt.Sprintf("\x1b[38;2;%d;%d;%dm", r, g, b)
			reset := "\x1b[0m"
			fmt.Println(ansiColor + line + reset)
		}
		color.Unset()
	},
	
	

}


func Execute() {
	
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}


