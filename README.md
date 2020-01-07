# cidr-to-glob

Convert CIDRs to glob-style patterns.

## Usage

```
cidr-to-glob [CIDR]...
cidr-to-glob -f file
command | cidr-to-glob
```

## Flags

```
-f string
      read from file instead of stdin
-o string
      write to file instead of stdout
```

## Example

```
# cidr-to-glob 10.0.0.0/9
10.[0-9].[0-9]*.[0-9]*
10.[1-9][0-9].[0-9]*.[0-9]*
10.1[0-1][0-9].[0-9]*.[0-9]*
10.12[0-7].[0-9]*.[0-9]*
```
