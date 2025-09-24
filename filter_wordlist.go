package main

import (
    "bufio"
    "fmt"
    "log"
    "os"
    "regexp"
)

func main() {
    if len(os.Args) != 2 {
        fmt.Fprintf(os.Stderr, "Usage: %s <wordlist>\n", os.Args[0])
        os.Exit(1)
    }
    infile := os.Args[1]
    f, err := os.Open(infile)
    if err != nil {
        log.Fatalf("opening input: %v", err)
    }
    defer f.Close()

    // Compile all patterns once
    patterns := []string{
        // Original noise filters
        `[\!(,%]`, `.{100,}`, `[0-9]{4,}`, `[0-9]{3,}$`,
        `[a-z0-9]{32}`, `[0-9]+[A-Z0-9]{5,}`, `/.*?/.*?/.*?/.*?/.*?/.*?/`,
        `\w{8}-\w{4}-\w{4}-\w{4}-\w{12}`,
        `[0-9]+[a-zA-Z]+[0-9]+[a-zA-Z]+[0-9]+`,
        `\.(png|jpg|jpeg|gif|svg|bmp|ttf|avif|wav|mp4|aac|ajax|css|all)$`,
        `^$`,
        // First supplemental set
        `[^a-z0-9\.-]`, `[A-Z]`, `(^-|-$)`, `^[0-9]+$`, `^.{64,}$`,
        `.*\..*\..*`, `^xn--`, `^https?://`, `:\d+$`, `[?=&]`, `@`,
        `^\*\.`,
        // Second deeper set
        `[^[:print:]]`, `_`, `[\[\]\(\)\{\}]`, `["'<>]`, `[;,$]`, `\+`,
        `\.\.`, `(^\.|\.$)`, `^-+$`, `-{3,}`, `^(com|net|org|io|app|dev|local)$`,
        `^.{1,3}$`, // awk will also enforce ≥4, but we prefilter
        // Note: Go's regexp supports backrefs but we dropped that one
        // Additional pro‑level filters
        `%[0-9A-Fa-f]{2}`, `\|`, `\s`,
        `^(example|test|testing|sample|localhost)$`,
        `^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+$`,
        `^(dev|stage|staging|uat|qa|beta|prod)$`,
    }

    regs := make([]*regexp.Regexp, len(patterns))
    for i, p := range patterns {
        regs[i] = regexp.MustCompile(p)
    }

    scanner := bufio.NewScanner(f)
    writer := bufio.NewWriter(os.Stdout)
    defer writer.Flush()

    for scanner.Scan() {
        line := scanner.Text()
        L := len(line)
        if L < 4 || L > 32 {
            continue
        }

        skip := false
        for _, re := range regs {
            if re.MatchString(line) {
                skip = true
                break
            }
        }
        if skip {
            continue
        }

        fmt.Fprintln(writer, line)
    }
    if err := scanner.Err(); err != nil {
        log.Fatalf("reading input: %v", err)
    }
}
