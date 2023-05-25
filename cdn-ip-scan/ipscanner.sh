#!/bin/bash
# 请使用root安装zmap 否则会出现奇怪的问题
# Replace 'cidr_list.txt' with the actual path to your CIDR list file
filename="tw-hinet-cidr.txt"

function check_zmap() {
    if command -v zmap >/dev/null 2>&1; then
        echo "ZMap is already installed."
    else
        echo -e "\e[31mWarning: ZMap is not installed!\e[0m"
        exit 1
    fi
}

function append_file() {
   source_file="$1"
   destination_file="$2"

   if [ -f "$source_file" ]; then
       while IFS= read -r line; do
           echo "$line" >> "$destination_file"
       done < "$source_file"

       echo "Content of '$source_file' copied line by line and appended to '$destination_file'."
   else
       echo "Error: Source file '$source_file' does not exist."
       exit 1
   fi
}


function scan_ip_port_by_zmap(){
    cidr="$1"
    zmap_command="zmap -B 100K -p 443 $cidr -o results.csv"
    #echo "Running zmap scan for CIDR: $cidr"
    echo "Command: $zmap_command"

    eval "$zmap_command"

    #echo "Zmap scan completed for CIDR: $cidr"
    #zmap -B 100K -p 443 211.72.0.0/16 -G 00:16:3e:3e:c0:80 -o results.csv
    append_file resultips.txt scanned-results.csv
}


function process_cidr_file() {
    local filename="$1"

    while IFS= read -r line; do
        echo "Scanning the CIDR: $line"
        start_time=$(date +%s)
        scan_ip_port_by_zmap $line

        end_time=$(date +%s)
        time_taken=$((end_time - start_time))
        echo "Current CIDR: $line finished, time taken: $time_taken seconds"
    done < "$filename"
}

function main(){
    check_zmap
    task_start_time=$(date +%s)
    process_cidr_file $filename


    task_end_time=$(date +%s)
    task_time_taken=$((task_end_time - task_start_time))
    echo "----------------------------------------------------------------"
    echo "Task finished, time taken: $task_time_taken seconds"
}

main
