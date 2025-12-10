/**
 * 客户端MS2数据处理工具
 * 用于解析MS2数据并计算指纹，减轻服务器负担
 */

/**
 * 解析MS2文本数据，提取各能量级别的质谱峰
 * 
 * @param {string} ms2Text - MS2数据文本
 * @param {number} maxPeaks - 最大峰数量限制，0表示无限制
 * @param {boolean} keepTopIntensity - 是否保留强度最高的峰（true）还是简单截断（false）
 * @returns {Object} 字典，键为能量级别（如'energy0'），值为(m/z, intensity)列表
 *                   对于简单的两列格式文件（如A9_165640.txt），会返回 {'energy0': [(mz1, intensity1), ...]}
 */
export function parseMS2Data(ms2Text, maxPeaks = 2000, keepTopIntensity = true) {
    if (!ms2Text || typeof ms2Text !== 'string') {
        return {};
    }

    const ms2Data = {};
    let currentEnergy = null;
    const lines = ms2Text.trim().split('\n');
    
    // 检查是否为简单的两列格式（没有energy标记）
    // 查看前几行是否包含'energy'关键词
    const hasEnergyMarkers = lines.slice(0, 10).some(line => {
        const trimmedLine = line.trim();
        return trimmedLine && trimmedLine.startsWith('energy');
    });

    if (!hasEnergyMarkers) {
        // 简单两列格式，如A9_165640.txt
        // 创建一个默认的能量级别
        currentEnergy = 'energy0';
        const allPeaks = [];
        
        for (const line of lines) {
            const trimmedLine = line.trim();
            if (!trimmedLine) {
                continue;
            }
            
            // 跳过注释行
            if (trimmedLine.startsWith('#')) {
                continue;
            }
            
            // 解析两列数据：m/z intensity
            const parts = trimmedLine.split(/\s+/);
            if (parts.length < 2) {
                continue;
            }
            
            try {
                const mz = parseFloat(parts[0]);
                const intensity = parseFloat(parts[1]);
                
                if (!isNaN(mz) && !isNaN(intensity)) {
                    allPeaks.push([mz, intensity]);
                }
            } catch (error) {
                // 跳过无法解析的行
                console.warn('无法解析MS2数据行:', trimmedLine, error);
                continue;
            }
        }
        
        console.log(`解析了 ${allPeaks.length} 个峰`);
        
        // 处理峰数据
        let processedPeaks = allPeaks;
        if (maxPeaks > 0 && allPeaks.length > maxPeaks) {
            if (keepTopIntensity) {
                // 按强度降序排序，保留强度最高的峰
                processedPeaks = allPeaks
                    .sort((a, b) => b[1] - a[1]) // 按强度降序排序
                    .slice(0, maxPeaks);
                console.log(`保留强度最高的 ${maxPeaks} 个峰（原 ${allPeaks.length} 个）`);
            } else {
                // 简单截断
                processedPeaks = allPeaks.slice(0, maxPeaks);
                console.log(`截断保留前 ${maxPeaks} 个峰（原 ${allPeaks.length} 个）`);
            }
        }
        
        ms2Data[currentEnergy] = processedPeaks;
    } else {
        // 原始格式：包含energy标记
        const energyPeaks = {};
        
        for (const line of lines) {
            const trimmedLine = line.trim();
            if ((!trimmedLine) && currentEnergy.startsWith('energy2')) {
                break;
            }
            if (!trimmedLine) {
                continue;
            }
            
            // 检查是否为能量级别标记
            if (trimmedLine.startsWith('energy')) {
                currentEnergy = trimmedLine;
                energyPeaks[currentEnergy] = [];
                continue;
            }
            
            // 跳过注释行
            if (trimmedLine.startsWith('#')) {
                continue;
            }
            
            // 解析质谱峰数据行
            // 格式: m/z intensity fragment_ids (additional_data...)
            const parts = trimmedLine.split(/\s+/);
            if (parts.length < 2) {
                continue;
            }
            
            try {
                const mz = parseFloat(parts[0]);
                const intensity = parseFloat(parts[1]);
                
                if (!isNaN(mz) && !isNaN(intensity) && currentEnergy) {
                    // 添加到当前能量级别
                    if (!energyPeaks[currentEnergy]) {
                        energyPeaks[currentEnergy] = [];
                    }
                    energyPeaks[currentEnergy].push([mz, intensity]);
                }
            } catch (error) {
                // 跳过无法解析的行
                console.warn('无法解析MS2数据行:', trimmedLine, error);
                continue;
            }
        }
        
        // 处理每个能量级别的峰数据
        for (const energyLevel in energyPeaks) {
            if (energyPeaks.hasOwnProperty(energyLevel)) {
                const peaks = energyPeaks[energyLevel];
                console.log(`能量级别 ${energyLevel}: 解析了 ${peaks.length} 个峰`);
                
                let processedPeaks = peaks;
                if (maxPeaks > 0 && peaks.length > maxPeaks) {
                    if (keepTopIntensity) {
                        // 按强度降序排序，保留强度最高的峰
                        processedPeaks = peaks
                            .sort((a, b) => b[1] - a[1]) // 按强度降序排序
                            .slice(0, maxPeaks);
                        console.log(`能量级别 ${energyLevel}: 保留强度最高的 ${maxPeaks} 个峰`);
                    } else {
                        // 简单截断
                        processedPeaks = peaks.slice(0, maxPeaks);
                        console.log(`能量级别 ${energyLevel}: 截断保留前 ${maxPeaks} 个峰`);
                    }
                }
                
                ms2Data[energyLevel] = processedPeaks;
            }
        }
    }
    console.log("峰解析数据: ",ms2Data)
    return ms2Data;
}

/**
 * 计算MS2指纹向量
 * 
 * @param {Object} ms2Data - 解析后的MS2数据
 * @param {number} binSize - m/z分箱大小（Da），默认1.0
 * @param {number} maxMz - 最大m/z值，默认1000.0
 * @returns {Array<number>} MS2指纹向量
 */
export function calculateMS2Fingerprint(ms2Data, binSize = 1.0, maxMz = 1000.0) {
    if (!ms2Data || Object.keys(ms2Data).length === 0) {
        return [];
    }
    
    // 创建分箱
    const numBins = Math.floor(maxMz / binSize) + 1;
    const fingerprint = new Array(numBins).fill(0);
    
    // 合并所有能量级别的峰
    const allPeaks = [];
    for (const energyLevel in ms2Data) {
        if (ms2Data.hasOwnProperty(energyLevel)) {
            allPeaks.push(...ms2Data[energyLevel]);
        }
    }
    
    // 将峰分配到分箱中
    for (const [mz, intensity] of allPeaks) {
        const binIdx = Math.floor(mz / binSize);
        if (binIdx < numBins) {
            // 使用强度平方根作为权重（常见于质谱相似度计算）
            fingerprint[binIdx] += Math.sqrt(intensity);
        }
    }
    
    // 归一化
    const norm = Math.sqrt(fingerprint.reduce((sum, val) => sum + val * val, 0));
    if (norm > 0) {
        return fingerprint.map(val => val / norm);
    }
    
    return fingerprint;
}

/**
 * 将浮点数指纹转换为二进制指纹
 * 
 * @param {Array<number>} fp - 浮点数指纹数组
 * @param {number} threshold - 阈值，大于该值的设为1，否则为0，默认0.01
 * @returns {Array<number>} 二进制指纹数组（0和1）
 */
export function floatFingerprintToBinary(fp, threshold = 0.01) {
    return fp.map(val => val > threshold ? 1 : 0);
}

/**
 * 将二进制指纹转换为base64编码
 * 
 * @param {Array<number>} binaryFp - 二进制指纹数组（0和1）
 * @returns {string} base64编码的字符串
 */
export function binaryFingerprintToBase64(binaryFp) {
    // 将二进制数组转换为字节
    // 每8位转换为一个字节
    // 确保长度是8的倍数
    const paddedLength = Math.ceil(binaryFp.length / 8) * 8;
    const padded = new Uint8Array(paddedLength);
    padded.set(binaryFp.map(bit => bit ? 1 : 0));
    
    // 转换为字节
    const bytes = new Uint8Array(paddedLength / 8);
    for (let i = 0; i < bytes.length; i++) {
        let byte = 0;
        for (let j = 0; j < 8; j++) {
            byte = (byte << 1) | (padded[i * 8 + j] || 0);
        }
        bytes[i] = byte;
    }
    
    // 转换为base64
    const binaryString = String.fromCharCode(...bytes);
    return btoa(binaryString);
}

/**
 * 计算MS2指纹并返回包含base64编码的二进制指纹的JSON
 * 
 * @param {Object} ms2Data - 解析后的MS2数据
 * @returns {Object} 包含峰列表和base64编码的二进制指纹的对象
 */
export function calculateMS2FingerprintWithBase64(ms2Data) {
    // 合并所有能量级别的峰
    const allPeaks = [];
    for (const energyLevel in ms2Data) {
        if (ms2Data.hasOwnProperty(energyLevel)) {
            allPeaks.push(...ms2Data[energyLevel]);
        }
    }
    
    // 计算浮点数指纹
    const fingerprintFloat = calculateMS2Fingerprint(ms2Data);
    
    // 转换为二进制指纹
    const fingerprintBinary = floatFingerprintToBinary(fingerprintFloat, 0.01);
    
    // 转换为base64编码
    const fingerprintBase64 = binaryFingerprintToBase64(fingerprintBinary);
    
    return {
        peaks: allPeaks,
        fingerprint_base64: fingerprintBase64
    };
}

/**
 * 计算modified cosine相似度（简化版本，用于客户端预筛选）
 * 
 * @param {Array<Array<number>>} peaks1 - 第一个质谱的峰列表，每个元素为[m/z, intensity]
 * @param {Array<Array<number>>} peaks2 - 第二个质谱的峰列表，每个元素为[m/z, intensity]
 * @param {number} tolerance - 质量容差（Da），默认0.5
 * @returns {number} 相似度分数（0-1之间）
 */
export function modifiedCosineSimilarity(peaks1, peaks2, tolerance = 0.5) {
    if (!peaks1 || !peaks2 || peaks1.length === 0 || peaks2.length === 0) {
        return 0.0;
    }
    
    // 按m/z排序
    const peaks1Sorted = [...peaks1].sort((a, b) => a[0] - b[0]);
    const peaks2Sorted = [...peaks2].sort((a, b) => a[0] - b[0]);
    
    // 初始化指针
    let i = 0, j = 0;
    let matchedProductSum = 0.0;
    
    // 滑动窗口匹配
    while (i < peaks1Sorted.length && j < peaks2Sorted.length) {
        const [mz1, int1] = peaks1Sorted[i];
        const [mz2, int2] = peaks2Sorted[j];
        
        // 计算质量差
        const deltaMz = Math.abs(mz1 - mz2);
        
        if (deltaMz <= tolerance) {
            // 匹配成功，计算强度乘积
            matchedProductSum += int1 * int2;
            
            // 两个指针都向前移动
            i++;
            j++;
        } else if (mz1 < mz2) {
            // peaks1的m/z较小，移动i
            i++;
        } else {
            // peaks2的m/z较小，移动j
            j++;
        }
    }
    
    // 计算分母
    const sum1 = peaks1Sorted.reduce((sum, [_, int]) => sum + int * int, 0);
    const sum2 = peaks2Sorted.reduce((sum, [_, int]) => sum + int * int, 0);
    
    if (sum1 === 0 || sum2 === 0) {
        return 0.0;
    }
    
    // 计算modified cosine相似度
    const similarity = matchedProductSum / (Math.sqrt(sum1) * Math.sqrt(sum2));
    
    return Math.max(0.0, Math.min(1.0, similarity));
}

/**
 * 计算指纹的cosine相似度（用于快速预筛选）
 * 
 * @param {Array<number>} fp1 - 第一个指纹向量
 * @param {Array<number>} fp2 - 第二个指纹向量
 * @returns {number} 相似度分数（0-1之间）
 */
export function cosineSimilarityFingerprint(fp1, fp2) {
    if (!fp1 || !fp2 || fp1.length === 0 || fp2.length === 0 || fp1.length !== fp2.length) {
        return 0.0;
    }
    
    // 计算点积
    let dotProduct = 0;
    for (let i = 0; i < fp1.length; i++) {
        dotProduct += fp1[i] * fp2[i];
    }
    
    // 计算范数
    const norm1 = Math.sqrt(fp1.reduce((sum, val) => sum + val * val, 0));
    const norm2 = Math.sqrt(fp2.reduce((sum, val) => sum + val * val, 0));
    
    if (norm1 === 0 || norm2 === 0) {
        return 0.0;
    }
    
    return dotProduct / (norm1 * norm2);
}

/**
 * 处理MS2文本并计算指纹（包含base64编码的二进制指纹）
 * 
 * @param {string} ms2Text - MS2数据文本
 * @param {number} maxPeaks - 最大峰数量限制，0表示无限制
 * @param {boolean} keepTopIntensity - 是否保留强度最高的峰（true）还是简单截断（false）
 * @returns {Object} 包含解析后的数据和指纹的对象
 */
export function processMS2Text(ms2Text, maxPeaks = 5000, keepTopIntensity = true) {
    try {
        // 解析MS2数据
        const ms2Data = parseMS2Data(ms2Text, maxPeaks, keepTopIntensity);
        
        if (Object.keys(ms2Data).length === 0) {
            throw new Error('无法解析MS2数据：数据为空或格式不正确');
        }
        
        // 计算指纹（包含base64编码的二进制指纹）
        const fingerprintData = calculateMS2FingerprintWithBase64(ms2Data);
        
        // 计算浮点数指纹（向后兼容）
        const fingerprintFloat = calculateMS2Fingerprint(ms2Data);
        
        // 合并所有能量级别的峰
        const allPeaks = fingerprintData.peaks;
        
        return {
            success: true,
            ms2Data,
            fingerprint: fingerprintFloat, // 向后兼容
            peaks: allPeaks,
            fingerprintJson: JSON.stringify({
                peaks: allPeaks,
                fingerprint_base64: fingerprintData.fingerprint_base64,
                timestamp: new Date().toISOString()
            })
        };
    } catch (error) {
        console.error('处理MS2数据时出错:', error);
        return {
            success: false,
            error: error.message,
            ms2Data: {},
            fingerprint: [],
            peaks: [],
            fingerprintJson: ''
        };
    }
}

/**
 * 从MS2文本直接计算指纹JSON（简化接口）
 * 
 * @param {string} ms2Text - MS2数据文本
 * @param {number} maxPeaks - 最大峰数量限制，0表示无限制
 * @param {boolean} keepTopIntensity - 是否保留强度最高的峰（true）还是简单截断（false）
 * @returns {string} 指纹JSON字符串
 */
export function calculateMS2FingerprintJson(ms2Text, maxPeaks = 5000, keepTopIntensity = true) {
    const result = processMS2Text(ms2Text, maxPeaks, keepTopIntensity);
    return result.fingerprintJson;
}

export default {
    parseMS2Data,
    calculateMS2Fingerprint,
    modifiedCosineSimilarity,
    cosineSimilarityFingerprint,
    processMS2Text,
    calculateMS2FingerprintJson
};
