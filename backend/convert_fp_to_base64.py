# Re-import necessary libraries and recreate the CSV with fp column
import pandas as pd
import random
from rdkit import Chem
from rdkit.Chem import AllChem
import base64
import numpy as np

# Recreate the synthetic dataset
random.seed(42)
n = 100
rows = []
item_types = ["Small molecule", "Peptide", "Polymer", "Natural product", "Intermediate"]
tags = ["antibacterial","antiviral","anticancer","analgesic","antiinflammatory","enzyme inhibitor","CNS-active","unknown"]
bioactivities = ["Active", "Inactive", "Moderate", "Potent", "Weak"]
smiles_list = [
    "CCO","CC(=O)O","c1ccccc1","CCN(CC)CC","C1CCCCC1","O=C(O)C(C)N","CCOC(=O)C","CC(C)C(=O)O",
    "CC(=O)NC","CCS","C=CC(=O)O","CC#N","CCOCC","C1=CC=CN=C1","CCC(=O)O","C(C(=O)O)N","CCOC","CCN","COC(=O)C","CNC(=O)C"
]

def random_formula():
    c = random.randint(1,20)
    h = random.randint(0,40)
    o = random.randint(0,5)
    n_ = random.randint(0,3)
    parts = []
    parts.append(f"C{c}" if c>1 else "C")
    if h>0: parts.append(f"H{h}" if h>1 else "H")
    if o>0: parts.append(f"O{o}" if o>1 else "O")
    if n_>0: parts.append(f"N{n_}" if n_>1 else "N")
    return "".join(parts)

def random_iupac(sm):
    prefixes = ["2-", "3-", "4-","N-",""]
    bases = ["methyl","ethyl","propyl","butyl","benzyl","phenyl","acetyl","carboxyl","amino","oxo"]
    return f"{random.choice(prefixes)}{random.choice(bases)} {random.choice(['acid','amide','ester','chloride','sulfate'])}".strip()

for i in range(1, n+1):
    id_ = f"CMP{i:04d}"
    source = random.choice(["PubChem","ChEMBL","In-house","Literature","ZINC"])
    item_name = f"Compound {i:03d}"
    item_type = random.choice(item_types)
    smiles = random.choice(smiles_list)
    iupac = random_iupac(smiles)
    description = f"Synthetic sample {i}. {random.choice(['Test compound','Lead-like','Fragment','Designed ligand','Screening hit'])}."
    cas = f"{random.randint(100,99999)}-{random.randint(10,99)}-{random.randint(0,9)}"
    tag = random.choice(tags)
    formula = random_formula()
    structure = "2D" if random.random()<0.8 else "3D"
    ms1 = round(random.uniform(50,800),4)
    ms2 = round(ms1 * random.uniform(0.2,0.9),4)
    bio = random.choice(bioactivities)
    rows.append({
        "ID": id_,
        "Source": source,
        "Item Name": item_name,
        "Item Type": item_type,
        "IUPAC name": iupac,
        "Description": description,
        "CAS number": cas,
        "Item Tag": tag,
        "Formula": formula,
        "Structure": structure,
        "MS1": ms1,
        "MS2": ms2,
        "Bioactivity": bio,
        "smiles": smiles
    })

df = pd.DataFrame(rows)

# Function to generate 2048-bit Morgan fingerprint as base64
def smiles_to_fp_base64(smiles, nBits=2048):
    mol = Chem.MolFromSmiles(smiles)
    if mol is None:
        return ""
    fp = AllChem.GetMorganFingerprintAsBitVect(mol, radius=2, nBits=nBits)
    return fp.ToBase64()

# Add fingerprint column
df['fp'] = df['smiles'].apply(lambda x: smiles_to_fp_base64(x, 2048))

# Save to CSV
out_path_fp = "./compounds_100_fp.csv"
df.to_csv(out_path_fp, index=False)