# Example YAML Tagging Ruleset

Here’s both a **YAML tagging structure** and an **IEC 81346-style rule set** for your instrument/equipment IDs.

---

## 📄 **1. YAML Tagging Structure**

This structure defines a hierarchy for one instrument, organized by *System → Equipment → Instrument*, allowing easy parsing and extensibility.

```yaml
Line4:
  Catalyst:
    SystemCode: CAT
    Equipment:
      PMP-0101:
        Type: Pump
        Motor:
          Tag: CAT-MTR-0101
          Description: 1 HP, 480V, 3PH Motor
        Drive:
          Tag: CAT-VFD-0101
          Description: 2/5 HP WEG or TECO Drive
        FlowMeter:
          Tag: CAT-FM-0101
          Type: Coriolis
          Description: 1/2" Endress+Hauser, analog output
      SCL-0102:
        Type: Scale
        Instrument:
          Tag: CAT-WT-0102
          Description: 2500 lb Pennsylvania analog scale
```

This approach:

* Encourages **readability**
* Supports **JSON export**
* Can integrate with CMDBs or SCADA config generators

---

## 🏷 **2. IEC 81346-style Rule Set**

Below is a simplified adaptation of IEC 81346 multi-aspect modeling. Tags are composed of **Function**, **Product**, and optionally **Location**.

### 💡 Core Structure

```xml
<Function>-<EquipmentType>-<InstanceNumber>
```

Or in long form:

```xml
F=<SystemCode>/P=<EquipmentType>-<Instance>
```

### 📐 Tagging Rules

| Aspect       | Definition                          | Tag Segment          | Example                |
| ------------ | ----------------------------------- | -------------------- | ---------------------- |
| `Function`   | Process/System name or abbreviation | `CAT`, `POL`         | `CAT` = Catalyst       |
| `Product`    | Equipment or instrument type        | `PMP`, `FM`, etc.    | `FM` = Flow Meter      |
| `Instance`   | Numeric ID within that group        | `0101`, `0202`       | `0101` = First in CAT  |
| `Redundancy` | Optional suffix: `-A`, `-1`, etc.   | e.g. `CAT-FM-0101-A` | For backup instruments |

### 📎 Full Example Tags

| Tag            | Meaning                                |
| -------------- | -------------------------------------- |
| `CAT-FM-0101`  | Flow Meter #1 on Catalyst pump         |
| `POL-VLV-0203` | 3-way valve on Polyol pump             |
| `ISO-WT-0402`  | Weight transmitter on Isocyanate scale |
| `HD-MTR-0601`  | Motor on Mixing Head                   |
| `BLT-VFD-0501` | Drive on Conveyor Belt                 |

---

## 🧰 Implementation Tips

* Store these in a database with lookup keys for:

  * Function code → full system name
  * Equipment type → role & symbol (for P\&ID)
  * Instance → numeric, no reuse
* Generate loop drawings from tag hierarchy
* Validate tags against regex:
  `^[A-Z]{3}-[A-Z]{3}-\d{4}(-[A-Z0-9])?$`

---

Let me know if you'd like a script to auto-generate these tags from a CSV input or want to expand the YAML into a full tagging spec for a project.
