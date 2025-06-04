#!/bin/bash

# Configuration bun / deno / node
PROCESSES=("82873" "83187" "83600")  # Change these to your processes
CSV_FILE="ram_monitor.csv"
INTERVAL=5  # seconds between checks

# Create CSV header if file doesn't exist
if [ ! -f "$CSV_FILE" ]; then
    echo "timestamp,${PROCESSES[0]}_mb,${PROCESSES[1]}_mb,${PROCESSES[2]}_mb" > "$CSV_FILE"
fi

# Function to get RAM usage in MB for a process
get_ram_usage() {
    local process_name="$1"
    # Use ps to get RSS (Resident Set Size) in KB, then convert to MB
    ps -A -o pid,rss,comm | grep -i "$process_name" | head -1 | awk '{print $2/1024}'
}

# Monitor loop
while true; do
    timestamp=$(date '+%Y-%m-%d %H:%M:%S')
    
    # Get RAM usage for each process
    ram1=$(get_ram_usage "${PROCESSES[0]}")
    ram2=$(get_ram_usage "${PROCESSES[1]}")
    ram3=$(get_ram_usage "${PROCESSES[2]}")
    
    # Handle cases where process might not be found
    ram1=${ram1:-0}
    ram2=${ram2:-0}
    ram3=${ram3:-0}
    
    # Append to CSV
    echo "$timestamp,$ram1,$ram2,$ram3" >> "$CSV_FILE"
    
    echo "$(date): Logged RAM usage - ${PROCESSES[0]}: ${ram1}MB, ${PROCESSES[1]}: ${ram2}MB, ${PROCESSES[2]}: ${ram3}MB"
    
    sleep $INTERVAL
done