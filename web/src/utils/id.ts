/**
 * ID 生成工具函数
 * 提供多种 ID 生成方法
 */

// ============================================================================
// UUID 生成
// ============================================================================

/**
 * 生成 UUID v4
 * @example
 * uuid() // '550e8400-e29b-41d4-a716-446655440000'
 */
export function uuid(): string {
  // 使用 crypto API 如果可用
  if (typeof crypto !== "undefined" && crypto.randomUUID) {
    return crypto.randomUUID();
  }

  // 回退实现
  return "xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx".replace(/[xy]/g, (c) => {
    const r = (Math.random() * 16) | 0;
    const v = c === "x" ? r : (r & 0x3) | 0x8;
    return v.toString(16);
  });
}

/**
 * 生成简短 UUID（移除连字符）
 * @example
 * shortUuid() // '550e8400e29b41d4a716446655440000'
 */
export function shortUuid(): string {
  return uuid().replace(/-/g, "");
}

/**
 * 验证 UUID 格式
 * @example
 * isValidUuid('550e8400-e29b-41d4-a716-446655440000') // true
 */
export function isValidUuid(id: string): boolean {
  const uuidRegex = /^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/i;
  return uuidRegex.test(id);
}

// ============================================================================
// NanoID 风格生成
// ============================================================================

const DEFAULT_ALPHABET = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz";
const URL_SAFE_ALPHABET = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-_";

/**
 * 生成随机字节数组
 */
function getRandomBytes(size: number): Uint8Array {
  if (typeof crypto !== "undefined" && crypto.getRandomValues) {
    return crypto.getRandomValues(new Uint8Array(size));
  }

  // 回退实现
  const bytes = new Uint8Array(size);
  for (let i = 0; i < size; i++) {
    bytes[i] = Math.floor(Math.random() * 256);
  }
  return bytes;
}

/**
 * 生成 NanoID 风格的 ID
 * @example
 * nanoid() // 'V1StGXR8_Z5jdHi6B-myT'
 * nanoid(10) // 'IRFa-VaY2b'
 */
export function nanoid(size: number = 21): string {
  const bytes = getRandomBytes(size);
  let id = "";

  for (let i = 0; i < size; i++) {
    id += URL_SAFE_ALPHABET[bytes[i] & 63];
  }

  return id;
}

/**
 * 使用自定义字母表生成 ID
 * @example
 * customId(10, '0123456789') // '8294756103'
 */
export function customId(size: number, alphabet: string): string {
  const _bytes = getRandomBytes(size);
  const mask = (2 << (Math.log(alphabet.length - 1) / Math.LN2)) - 1;
  const step = Math.ceil((1.6 * mask * size) / alphabet.length);

  let id = "";

  while (id.length < size) {
    const randomBytes = getRandomBytes(step);
    for (let i = 0; i < step && id.length < size; i++) {
      const byte = randomBytes[i] & mask;
      if (byte < alphabet.length) {
        id += alphabet[byte];
      }
    }
  }

  return id;
}

/**
 * 生成字母数字 ID
 * @example
 * alphanumericId(8) // 'a1B2c3D4'
 */
export function alphanumericId(size: number = 12): string {
  return customId(size, DEFAULT_ALPHABET);
}

/**
 * 生成纯数字 ID
 * @example
 * numericId(6) // '123456'
 */
export function numericId(size: number = 6): string {
  return customId(size, "0123456789");
}

/**
 * 生成纯字母 ID（小写）
 * @example
 * alphabeticId(8) // 'abcdefgh'
 */
export function alphabeticId(size: number = 8): string {
  return customId(size, "abcdefghijklmnopqrstuvwxyz");
}

/**
 * 生成十六进制 ID
 * @example
 * hexId(16) // 'a1b2c3d4e5f67890'
 */
export function hexId(size: number = 16): string {
  return customId(size, "0123456789abcdef");
}

// ============================================================================
// 时间戳 ID
// ============================================================================

/**
 * 生成时间戳 ID
 * @example
 * timestampId() // '1701234567890_a1b2c3d4'
 */
export function timestampId(randomLength: number = 8): string {
  const timestamp = Date.now();
  const random = alphanumericId(randomLength);
  return `${timestamp}_${random}`;
}

/**
 * 生成有序 ID（时间戳 + 随机）
 * 适合需要按时间排序的场景
 * @example
 * sortableId() // '0h5k7g2a1b3c4d5e'
 */
export function sortableId(randomLength: number = 8): string {
  // 使用 base36 编码时间戳以节省空间
  const timestamp = Date.now().toString(36);
  const random = alphanumericId(randomLength);
  return `${timestamp}${random}`;
}

/**
 * 从有序 ID 提取时间戳
 * @example
 * extractTimestamp('lxyz1234abcd5678') // 1701234567890
 */
export function extractTimestamp(sortableId: string): number | null {
  // 尝试提取 base36 时间戳（前 8-9 个字符）
  const timestampPart = sortableId.slice(0, 9);
  const timestamp = parseInt(timestampPart, 36);

  if (isNaN(timestamp) || timestamp < 0) {
    return null;
  }

  return timestamp;
}

// ============================================================================
// 前缀 ID
// ============================================================================

/**
 * 生成带前缀的 ID
 * @example
 * prefixedId('user') // 'user_a1b2c3d4e5f6'
 * prefixedId('order', 16) // 'order_a1b2c3d4e5f67890'
 */
export function prefixedId(prefix: string, size: number = 12): string {
  return `${prefix}_${alphanumericId(size)}`;
}

/**
 * 创建带前缀的 ID 生成器
 * @example
 * const userIdGen = createPrefixedIdGenerator('user')
 * userIdGen() // 'user_a1b2c3d4e5f6'
 */
export function createPrefixedIdGenerator(prefix: string, size: number = 12): () => string {
  return () => prefixedId(prefix, size);
}

// ============================================================================
// 序列 ID
// ============================================================================

/**
 * 创建序列 ID 生成器
 * @example
 * const seq = createSequence()
 * seq() // 1
 * seq() // 2
 * seq() // 3
 */
export function createSequence(start: number = 1): () => number {
  let current = start;
  return () => current++;
}

/**
 * 创建带格式的序列 ID 生成器
 * @example
 * const orderSeq = createFormattedSequence('ORD', 6)
 * orderSeq() // 'ORD000001'
 * orderSeq() // 'ORD000002'
 */
export function createFormattedSequence(prefix: string, padding: number = 6, start: number = 1): () => string {
  let current = start;
  return () => {
    const id = String(current++).padStart(padding, "0");
    return `${prefix}${id}`;
  };
}

// ============================================================================
// 雪花 ID（简化版）
// ============================================================================

/**
 * 雪花 ID 生成器配置
 */
export interface SnowflakeConfig {
  /** 机器 ID (0-1023) */
  machineId?: number;
  /** 数据中心 ID (0-31) */
  datacenterId?: number;
  /** 自定义纪元 (默认: 2024-01-01) */
  epoch?: number;
}

/**
 * 创建雪花 ID 生成器
 * 生成趋势递增的 64 位 ID
 * @example
 * const snowflake = createSnowflake({ machineId: 1 })
 * snowflake() // '7159558526853120001'
 */
export function createSnowflake(config: SnowflakeConfig = {}): () => string {
  const {
    machineId = 1,
    datacenterId = 1,
    epoch = 1704067200000, // 2024-01-01
  } = config;

  let lastTimestamp = -1;
  let sequence = 0;

  const machineIdBits = 5;
  const datacenterIdBits = 5;
  const sequenceBits = 12;

  const maxMachineId = (1 << machineIdBits) - 1;
  const maxDatacenterId = (1 << datacenterIdBits) - 1;
  const maxSequence = (1 << sequenceBits) - 1;

  if (machineId < 0 || machineId > maxMachineId) {
    throw new Error(`Machine ID must be between 0 and ${maxMachineId}`);
  }
  if (datacenterId < 0 || datacenterId > maxDatacenterId) {
    throw new Error(`Datacenter ID must be between 0 and ${maxDatacenterId}`);
  }

  return () => {
    let timestamp = Date.now() - epoch;

    if (timestamp === lastTimestamp) {
      sequence = (sequence + 1) & maxSequence;
      if (sequence === 0) {
        // 等待下一毫秒
        while (timestamp <= lastTimestamp) {
          timestamp = Date.now() - epoch;
        }
      }
    } else {
      sequence = 0;
    }

    lastTimestamp = timestamp;

    // 使用 BigInt 进行位运算以避免精度丢失
    const id =
      (BigInt(timestamp) << BigInt(sequenceBits + machineIdBits + datacenterIdBits)) |
      (BigInt(datacenterId) << BigInt(sequenceBits + machineIdBits)) |
      (BigInt(machineId) << BigInt(sequenceBits)) |
      BigInt(sequence);

    return id.toString();
  };
}

// ============================================================================
// ULID（简化实现）
// ============================================================================

const ULID_ALPHABET = "0123456789ABCDEFGHJKMNPQRSTVWXYZ";

/**
 * 生成 ULID
 * 可排序的唯一 ID，26 字符
 * @example
 * ulid() // '01ARZ3NDEKTSV4RRFFQ69G5FAV'
 */
export function ulid(): string {
  const timestamp = Date.now();

  // 编码时间戳部分（10 字符）
  let timestampPart = "";
  let t = timestamp;
  for (let i = 0; i < 10; i++) {
    timestampPart = ULID_ALPHABET[t % 32] + timestampPart;
    t = Math.floor(t / 32);
  }

  // 生成随机部分（16 字符）
  const bytes = getRandomBytes(10);
  let randomPart = "";
  for (let i = 0; i < 16; i++) {
    const byteIndex = Math.floor(i / 2);
    const nibble = i % 2 === 0 ? bytes[byteIndex] >> 4 : bytes[byteIndex] & 0x0f;
    randomPart += ULID_ALPHABET[nibble];
  }

  return timestampPart + randomPart;
}

/**
 * 从 ULID 提取时间戳
 * @example
 * extractUlidTimestamp('01ARZ3NDEKTSV4RRFFQ69G5FAV') // 1469918176385
 */
export function extractUlidTimestamp(ulidStr: string): number {
  if (ulidStr.length !== 26) {
    throw new Error("Invalid ULID length");
  }

  const timestampPart = ulidStr.slice(0, 10);
  let timestamp = 0;

  for (let i = 0; i < 10; i++) {
    const char = timestampPart[i];
    const value = ULID_ALPHABET.indexOf(char);
    if (value === -1) {
      throw new Error(`Invalid ULID character: ${char}`);
    }
    timestamp = timestamp * 32 + value;
  }

  return timestamp;
}

// ============================================================================
// 实用工具
// ============================================================================

/**
 * 生成唯一的 DOM 元素 ID
 * @example
 * uniqueDomId('input') // 'input-a1b2c3d4'
 */
export function uniqueDomId(prefix: string = "el"): string {
  return `${prefix}-${alphanumericId(8).toLowerCase()}`;
}

/**
 * 创建 ID 生成器工厂
 * @example
 * const idFactory = createIdFactory({
 *   user: () => prefixedId('usr'),
 *   order: () => prefixedId('ord'),
 *   session: () => nanoid()
 * })
 * idFactory.user() // 'usr_a1b2c3d4e5f6'
 */
export function createIdFactory<T extends Record<string, () => string>>(generators: T): T {
  return generators;
}

/**
 * 批量生成 ID
 * @example
 * generateIds(5) // ['a1b2...', 'c3d4...', ...]
 * generateIds(3, () => uuid()) // [uuid(), uuid(), uuid()]
 */
export function generateIds(count: number, generator: () => string = nanoid): string[] {
  return Array.from({ length: count }, generator);
}

/**
 * 检查 ID 是否唯一（在给定集合中）
 * @example
 * isUniqueId('abc123', new Set(['def456', 'ghi789'])) // true
 */
export function isUniqueId(id: string, existingIds: Set<string>): boolean {
  return !existingIds.has(id);
}

/**
 * 生成唯一 ID（确保在集合中唯一）
 * @example
 * const existingIds = new Set(['abc'])
 * ensureUniqueId(existingIds) // 生成一个不在集合中的 ID
 */
export function ensureUniqueId(
  existingIds: Set<string>,
  generator: () => string = nanoid,
  maxAttempts: number = 100
): string {
  for (let i = 0; i < maxAttempts; i++) {
    const id = generator();
    if (!existingIds.has(id)) {
      return id;
    }
  }
  throw new Error(`Failed to generate unique ID after ${maxAttempts} attempts`);
}
