const c = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_";

function writeBits(value: number): number[] {
    let e = value;
    let t = 0;
    const bits: number[] = [];

    for (let a = 1; a <= 5; a++) {
        const n = 5 * a;
        if (e < (1 << n)) {
            let nArray = [];
            while (nArray.length < a) {
                let t = e & 63;
                e >>= 6;
                nArray.unshift(t);
            }
            nArray[0] = nArray[0] | t;
            bits.push(...nArray);
            return bits;
        }
        e -= 1 << n;
        t = (64 | t) >> 1;
    }
    throw new Error("Value too large to encode.");
}

function writeTalents(talentStr: string): number[] {
    const bits: number[] = [];
    const trees = talentStr.split('-');

    for (let a = 0; a < 3; a++) {
        const tree = trees[a] || '';
        bits.push(...writeBits(tree.length));

        let l = 0;
        while (l < tree.length) {
            let chunk = 0;
            let s = 0;
            while (s < 7 && l < tree.length) {
                const digit = parseInt(tree[l], 10);
                chunk = (chunk << 3) | digit;
                l++;
                s++;
            }
            bits.push(...writeBits(chunk));
        }
    }

    return bits;
}

// Function to write glyphs (reverse of parseGlyphs)
function writeGlyphs(glyphIds: number[]): string {
    const base32 = '0123456789abcdefghjkmnpqrstvwxyz'; // Base32 character set
    let glyphStr = '0'; // insert random 0

    for (let i = 0; i < glyphIds.length; i++) {
        const spellId = glyphIds[i];
        if (spellId) {
            const glyphSlotChar = base32[i];
            const c1 = (spellId >> 15) & 31;
            const c2 = (spellId >> 10) & 31;
            const c3 = (spellId >> 5) & 31;
            const c4 = spellId & 31;

            if (c1 < 0 || c2 < 0 || c3 < 0 || c4 < 0) {
                continue; // Invalid spell ID
            }

            glyphStr += glyphSlotChar +
                base32[c1] +
                base32[c2] +
                base32[c3] +
                base32[c4];
        }
    }

    return glyphStr;
}

// Function to write the hash (reverse of readHash)
function writeHash(data: any): string {
    let hash = '';

    // Initialize bits array
    const bits: number[] = [];

    // Starting character (B for gear planner)
    let idx = 1; // Assuming idx is 1 for gear planner

    // Include idx in the hash (as first character)
    hash += c[idx];

    // Gender (assuming genderId is 0 or 1)
    bits.push(0);

    // Level
    bits.push(...writeBits(data.level ?? 0));

    // Talents
    const talentBits = writeTalents(data.talents.join('-'));
    bits.push(...talentBits);

    // Glyphs
    const glyphStr = writeGlyphs(data.glyphs ?? []);
    const glyphBytes = glyphStr.split('').map(ch => c.indexOf(ch));
    bits.push(...writeBits(glyphBytes.length));
    bits.push(...glyphBytes);

    // Items
    const items = data.items ?? [];
    bits.push(...writeBits(items.length));

    for (const item of items) {
        let e = 0;
        const itemBits: number[] = [];

        // Encode flags into e
        if (item.randomEnchantId) e |= 1 << 6;
        if (item.reforge) e |= 1 << 5;
        const gemCount = Math.min((item.gemItemIds ?? []).length, 7);
        e |= gemCount << 2;
        const enchantCount = Math.min((item.enchantIds ?? []).length, 3);
        e |= enchantCount;

        // Item slot and ID
        itemBits.push(...writeBits(item.slotId ?? 0));
        itemBits.push(...writeBits(item.itemId ?? 0));

        // Random Enchant ID
        if (item.randomEnchantId) {
            let enchant = item.randomEnchantId;
            let negative = enchant < 0 ? 1 : 0;
            if (negative) enchant *= -1;
            enchant = (enchant << 1) | negative;
            itemBits.push(...writeBits(enchant));
        }

        // Reforge
        if (item.reforge) {
            itemBits.push(...writeBits(item.reforge));
        }

        // Gems
        const gems = item.gemItemIds ?? [];
        for (let i = 0; i < gemCount; i++) {
            itemBits.push(...writeBits(gems[i]));
        }

        // Enchants
        const enchants = item.enchantIds ?? [];
        for (let i = 0; i < enchantCount; i++) {
            itemBits.push(...writeBits(enchants[i]));
        }

        // e is the item flags; add it at the start of itemBits
        bits.push(...writeBits(e));
        bits.push(...itemBits);
    }

    // Encode bits into characters
    let hashData = '';
    for (const bit of bits) {
        hashData += c.charAt(bit);
    }

    // Append the hash data to the URL
    if (hashData) {
        hash += hashData;
    }

    return hash;
}



// Taken from Wowhead
function readBits(e: number[]): number {
    if (!e.length) return 0;
    let t = 0,
        a = 1,
        n = e[0];
    while ((32 & n) > 0) {
        a++;
        n <<= 1;
    }
    let l = 63 >> a,
        s = e.shift()! & l;
    a--;
    for (let n = 1; n <= a; n++) {
        t += 1 << (5 * n);
        s = (s << 6) | (e.shift() || 0);
    }
    return s + t;
}

// Taken from Wowhead
function parseTalents(e: number[]): string {
    let t = "";
    for (let a = 0; a < 3; a++) {
        let len = readBits(e);
        while (len > 0) {
            let n = "",
                l = readBits(e);
            while (len > 0 && n.length < 7) {
                let digit = 7 & l;
                l >>= 3;
                n = `${digit}${n}`;
                len--;
            }
            t += n;
        }
        t += "-";
    }
    return t.replace(/-+$/, "");
}

// Taken from Wowhead
function readHash(e: string): any {
    const t: any = {};
    const a = e.match(/^#?~(-?\d+)$/);
    if (a) {
        return t;
    }
    let n = /^([a-z-]+)\/([a-z-]+)(?:\/([a-zA-Z0-9_-]+))?$/.exec(e);
    if (!n) return t;
    {
        t.class = n[1];
    }
    {
        t.race = n[2];
    }
    let o = n[3];
    if (!o) return t;
    let idx = c.indexOf(o.substring(0, 1));
    o = o.substring(1);
    if (!o.length) return t;
    let u: number[] = [];
    for (let e = 0; e < o.length; e++)
        u.push(c.indexOf(o.substring(e, e + 1)));
    if (idx > 1) return t;
    {
        let e = readBits(u) - 1;
        if (e >= 0) t.genderId = e;
    }
    {
        let e = readBits(u);
        if (e) t.level = e;
    }
    {
        let e = parseTalents(u),
            a = readBits(u),
            n = u.splice(0, a).map(e => c[e]).join(""),
            l = e + (a ? `_${n}` : "");
        if ("" !== l) t.talentHash = l;
        if(t.talentHash) {
            let talents = parseTalentString(t.talentHash);
            t.talents = talents.talents;
            t.glyphs = talents.glyphs;
        }

    }
    {
        let itemCount = readBits(u);
        t.items = [];
        while (itemCount--) {
            let e: number,
                a: any = {};
            if (idx < 1) {
                e = u.shift()!;
                let t = (e >> 5) & 1;
                e &= 31;
                if (t) e |= 64;
            } else e = readBits(u);
            a.slotId = readBits(u);
            a.itemId = readBits(u);
            if (0 != ((e >> 6) & 1)) {
                let enchant = readBits(u),
                    t = 1 & enchant;
                enchant >>= 1;
                if (t) enchant *= -1;
                a.randomEnchantId = enchant;
            }
            if (0 != ((e >> 5) & 1)) a.reforge = readBits(u);
            {
                let gemCount = (e >> 2) & 7;
                while (gemCount--) {
                    if (!a.gemItemIds) a.gemItemIds = [];
                    a.gemItemIds.push(readBits(u));
                }
            }
            {
                let enchantCount = e & 3;
                while (enchantCount--) {
                    if (!a.enchantIds) a.enchantIds = [];
                    a.enchantIds.push(readBits(u));
                }
            }
            t.items.push(a);
        }
    }
    return t;
}

// Function to parse glyphs from the glyph string
function parseGlyphs(glyphStr: string): number[] {
    const glyphIds = Array(9).fill(0); // Nine potential glyph slots
    const base32 = '0123456789abcdefghjkmnpqrstvwxyz'; // Base32 character set
    let cur = 1; // we skip the first index for whatever reason

    while (cur < glyphStr.length) {
        // Get glyph slot index
        const glyphSlotChar = glyphStr[cur];
        const glyphSlotIndex = base32.indexOf(glyphSlotChar);
        cur++;

        if (glyphSlotIndex < 0 || glyphSlotIndex >= glyphIds.length) {
            continue; // Skip invalid glyph slots
        }

        if (cur + 4 > glyphStr.length) {
            break; // Not enough characters for a glyph ID
        }

        // Decode the spellId using base32 encoding (each character represents 5 bits)
        const c1 = base32.indexOf(glyphStr[cur]);
        const c2 = base32.indexOf(glyphStr[cur + 1]);
        const c3 = base32.indexOf(glyphStr[cur + 2]);
        const c4 = base32.indexOf(glyphStr[cur + 3]);
        cur += 4;

        if (c1 < 0 || c2 < 0 || c3 < 0 || c4 < 0) {
            continue; // Invalid character in spell ID
        }

        const spellId = (c1 << 15) | (c2 << 10) | (c3 << 5) | c4;

        glyphIds[glyphSlotIndex] = spellId;
    }

    return glyphIds;
}

function parseTalentString(talentString: string): { talents: string, glyphs: number[] } {
    const [talentPart, glyphPart] = talentString.split('_');

    // Parse the talents
    // Talent string is something like '001-2301-33223203120220120321'
    // Each part separated by '-' corresponds to a talent tree
    const talents = talentPart;

    // Parse the glyphs
    let glyphs: number[] = [];
    if (glyphPart) {
        glyphs = parseGlyphs(glyphPart);
    }

    return { talents, glyphs };
}

export interface WowheadGearPlannerData {
    class?: string;          
    race?: string;           
    genderId?: number;        
    level: number;            
    talents: string[];       
    glyphs: number[];        
    items: WowheadItemData[]; 
}

export interface WowheadItemData {
    slotId: number;         
    itemId: number;          
    randomEnchantId?: number; 
    reforge?: number;   
    gemItemIds?: number[];
    enchantIds?: number[]; 
}

export function parseWowheadGearLink(link: string): any {
    // Extract the part after 'cata/gear-planner/'
    const match = link.match(/cata\/gear-planner\/(.+)/);
    if (!match) {
        throw new Error(`Invalid WCL URL ${link}, must look like "https://www.wowhead.com/cata/gear-planner/CLASS/RACE/XXXX"`);
    }
    let e = match[1];
    return readHash(e);
}

export function createWowheadGearPlannerLink(data: WowheadGearPlannerData): string {
    const baseUrl = '';
    const hash = writeHash(data);
    return baseUrl + hash;
}