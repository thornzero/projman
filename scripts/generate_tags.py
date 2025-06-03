#!/usr/bin/env python3

import csv
import yaml
import argparse
from collections import defaultdict

def generate_tags(input_csv, output_yaml):
    instance_counters = {}
    yaml_structure = {}

    instrument_table = []

    with open(input_csv, newline='') as csvfile:
        reader = csv.DictReader(csvfile)
        for row in reader:
            system = row['System']
            system_code = row['SystemCode']
            equip_type = row['EquipType']
            function = row['Function']
            instrument = row['Instrument']
            description = row['Description']

            key = f"{system_code}-{equip_type}"
            instance_num = instance_counters.get(key, 101)
            equip_id = f"{system_code}-{equip_type}-{instance_num:04d}"
            instr_id = f"{system_code}-{instrument}-{instance_num:04d}"
            instance_counters[key] = instance_num + 1

            # YAML structure
            sys_branch = yaml_structure.setdefault('Line4', {}).setdefault(system, {
                'SystemCode': system_code,
                'Equipment': {}
            })
            equip_branch = sys_branch['Equipment'].setdefault(f"{equip_type}-{instance_num:04d}", {
                'Type': function
            })
            equip_branch[instrument] = {
                'Tag': instr_id,
                'Description': description
            }

            instrument_table.append({
                'System': system,
                'Equipment ID': equip_id,
                'Instrument ID': instr_id,
                'Description': description
            })

    with open(output_yaml, 'w') as yamlfile:
        yaml.dump(yaml_structure, yamlfile, sort_keys=False)

    print(f"Tagging complete. YAML written to {output_yaml}")

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Generate instrument/equipment tags from a CSV file.")
    parser.add_argument("input_csv", help="Input CSV file path")
    parser.add_argument("output_yaml", help="Output YAML file path")
    args = parser.parse_args()

    generate_tags(args.input_csv, args.output_yaml)
