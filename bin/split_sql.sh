
#!/bin/bash

# Name of the input file
input_file="daily_closes.sql"

# Check if the input file exists
if [ ! -f "$input_file" ]; then
    echo "Error: $input_file not found!"
    exit 1
fi

# Counter for INSERT statements
count=0

# Counter for output files
file_number=1

# Name of the current output file
output_file="daily_closes_${file_number}.sql"

# Read the input file line by line
while IFS= read -r line
do
    # Write the line to the current output file
    echo "$line" >> "$output_file"
    
    # If the line starts with "INSERT INTO", increment the counter
    if [[ $line == INSERT\ INTO* ]]; then
        ((count++))
    fi
    
    # If we've reached 10 INSERT statements, start a new file
    if [ $count -eq 10 ]; then
        count=0
        ((file_number++))
        output_file="daily_closes_${file_number}.sql"
    fi
done < "$input_file"

echo "Splitting complete. Check the daily_closes_*.sql files."
