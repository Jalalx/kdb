#/bin/bash

while IFS= read -r line; do
    kdb --embed "$line"
done < ./gk_training_data.txt