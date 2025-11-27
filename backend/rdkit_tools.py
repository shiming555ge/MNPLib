from rdkit import Chem, DataStructs
from rdkit.Chem import AllChem

# smiles式转存pdb
def smiles_to_pdb(smiles, out_file="output.pdb"):
    mol = Chem.MolFromSmiles(smiles)
    mol = Chem.AddHs(mol)                       # 加氢
    AllChem.EmbedMolecule(mol)                 # 3D 构象
    AllChem.UFFOptimizeMolecule(mol)           # UFF 优化
    Chem.MolToPDBFile(mol, out_file)
    return out_file

# 分子指纹相似度
def tanimoto_similarity(smiles1, smiles2):
    mol1 = Chem.MolFromSmiles(smiles1)
    mol2 = Chem.MolFromSmiles(smiles2)

    fp1 = AllChem.GetMorganFingerprintAsBitVect(mol1, radius=2, nBits=2048)
    fp2 = AllChem.GetMorganFingerprintAsBitVect(mol2, radius=2, nBits=2048)

    return DataStructs.TanimotoSimilarity(fp1, fp2)

# 判断子结构
def is_substructure(smarts_pattern, smiles):
    patt = Chem.MolFromSmarts(smarts_pattern)
    mol = Chem.MolFromSmiles(smiles)
    return mol.HasSubstructMatch(patt)

# 批量相似度搜索
def similarity_search(query_smiles, library, threshold=0.5):
    qmol = Chem.MolFromSmiles(query_smiles)
    qfp  = AllChem.GetMorganFingerprintAsBitVect(qmol, 2, 2048)
    
    results = []
    for s in library:
        mol = Chem.MolFromSmiles(s)
        fp = AllChem.GetMorganFingerprintAsBitVect(mol, 2, 2048)
        sim = DataStructs.TanimotoSimilarity(qfp, fp)
        if sim >= threshold:
            results.append((s, sim))
    return sorted(results, key=lambda x: -x[1])

# 批量子结构搜索
def substructure_search(pattern_smarts, library):
    patt = Chem.MolFromSmarts(pattern_smarts)
    result = []
    for s in library:
        mol = Chem.MolFromSmiles(s)
        if mol.HasSubstructMatch(patt):
            result.append(s)
    return result

if __name__=="__main__":
    # 模式
    while(cmd := input()):
        if(cmd == "SimSearch"):
            data=input("waiting...")
        