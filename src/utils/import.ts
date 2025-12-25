/**
 * CSV 导入工具
 * 用于解析 CSV 文件并验证数据
 */

export interface ParsedUser {
  username: string;
  email: string;
  password: string;
  full_name?: string;
  status?: string;
}

export interface ParseResult<T> {
  data: T[];
  errors: ParseError[];
  totalRows: number;
  validRows: number;
}

export interface ParseError {
  row: number;
  field: string;
  message: string;
}

/**
 * 解析 CSV 文件内容
 * @param content CSV 文件内容字符串
 * @returns 解析后的数据行数组
 */
export function parseCSV(content: string): string[][] {
  const lines = content.trim().split(/\r?\n/);
  return lines.map((line) => {
    const result: string[] = [];
    let current = "";
    let inQuotes = false;

    for (let i = 0; i < line.length; i++) {
      const char = line[i];
      const nextChar = line[i + 1];

      if (inQuotes) {
        if (char === '"' && nextChar === '"') {
          current += '"';
          i++;
        } else if (char === '"') {
          inQuotes = false;
        } else {
          current += char;
        }
      } else {
        if (char === '"') {
          inQuotes = true;
        } else if (char === ",") {
          result.push(current.trim());
          current = "";
        } else {
          current += char;
        }
      }
    }
    result.push(current.trim());
    return result;
  });
}

/**
 * 验证邮箱格式
 */
function isValidEmail(email: string): boolean {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  return emailRegex.test(email);
}

/**
 * 验证用户名格式
 */
function isValidUsername(username: string): boolean {
  // 用户名: 3-50字符，只允许字母、数字、下划线
  const usernameRegex = /^[a-zA-Z0-9_]{3,50}$/;
  return usernameRegex.test(username);
}

/**
 * 验证密码强度
 */
function isValidPassword(password: string): boolean {
  // 密码: 至少8字符，包含大小写字母和数字
  if (password.length < 8) return false;
  const hasUpperCase = /[A-Z]/.test(password);
  const hasLowerCase = /[a-z]/.test(password);
  const hasNumber = /[0-9]/.test(password);
  return hasUpperCase && hasLowerCase && hasNumber;
}

/**
 * 解析并验证用户 CSV 数据
 * @param content CSV 文件内容
 * @returns 解析结果，包含有效数据和错误信息
 */
export function parseUserCSV(content: string): ParseResult<ParsedUser> {
  const rows = parseCSV(content);
  const errors: ParseError[] = [];
  const data: ParsedUser[] = [];

  if (rows.length === 0) {
    return { data: [], errors: [{ row: 0, field: "", message: "CSV 文件为空" }], totalRows: 0, validRows: 0 };
  }

  // 解析表头
  const headerRow = rows[0];
  if (!headerRow) {
    return { data: [], errors: [{ row: 0, field: "", message: "CSV 文件为空" }], totalRows: 0, validRows: 0 };
  }
  const headers = headerRow.map((h) => h.toLowerCase().trim());
  const usernameIdx = headers.indexOf("username");
  const emailIdx = headers.indexOf("email");
  const passwordIdx = headers.indexOf("password");
  const fullNameIdx = headers.indexOf("full_name");
  const statusIdx = headers.indexOf("status");

  // 验证必需列
  if (usernameIdx === -1) {
    errors.push({ row: 1, field: "header", message: "缺少必需列: username" });
  }
  if (emailIdx === -1) {
    errors.push({ row: 1, field: "header", message: "缺少必需列: email" });
  }
  if (passwordIdx === -1) {
    errors.push({ row: 1, field: "header", message: "缺少必需列: password" });
  }

  if (errors.length > 0) {
    return { data: [], errors, totalRows: rows.length - 1, validRows: 0 };
  }

  // 解析数据行
  for (let i = 1; i < rows.length; i++) {
    const row = rows[i];
    if (!row) continue;

    const rowNum = i + 1;
    let hasError = false;

    // 跳过空行
    if (row.length === 1 && row[0] === "") continue;

    const username = row[usernameIdx] || "";
    const email = row[emailIdx] || "";
    const password = row[passwordIdx] || "";
    const fullName = fullNameIdx !== -1 ? row[fullNameIdx] || "" : "";
    const status = statusIdx !== -1 ? row[statusIdx] || "active" : "active";

    // 验证 username
    if (!username) {
      errors.push({ row: rowNum, field: "username", message: "用户名不能为空" });
      hasError = true;
    } else if (!isValidUsername(username)) {
      errors.push({ row: rowNum, field: "username", message: "用户名格式无效（3-50字符，仅字母数字下划线）" });
      hasError = true;
    }

    // 验证 email
    if (!email) {
      errors.push({ row: rowNum, field: "email", message: "邮箱不能为空" });
      hasError = true;
    } else if (!isValidEmail(email)) {
      errors.push({ row: rowNum, field: "email", message: "邮箱格式无效" });
      hasError = true;
    }

    // 验证 password
    if (!password) {
      errors.push({ row: rowNum, field: "password", message: "密码不能为空" });
      hasError = true;
    } else if (!isValidPassword(password)) {
      errors.push({ row: rowNum, field: "password", message: "密码需至少8字符，包含大小写字母和数字" });
      hasError = true;
    }

    // 验证 status
    const validStatuses = ["active", "inactive", "suspended"];
    if (status && !validStatuses.includes(status.toLowerCase())) {
      errors.push({ row: rowNum, field: "status", message: `状态无效，有效值: ${validStatuses.join(", ")}` });
      hasError = true;
    }

    if (!hasError) {
      data.push({
        username,
        email,
        password,
        full_name: fullName || undefined,
        status: status.toLowerCase() || "active",
      });
    }
  }

  return {
    data,
    errors,
    totalRows: rows.length - 1,
    validRows: data.length,
  };
}

/**
 * 读取文件为文本
 * @param file File 对象
 * @returns Promise<string> 文件内容
 */
export function readFileAsText(file: File): Promise<string> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onload = (e) => {
      const result = e.target?.result;
      if (typeof result === "string") {
        resolve(result);
      } else {
        reject(new Error("无法读取文件内容"));
      }
    };
    reader.onerror = () => reject(new Error("文件读取失败"));
    reader.readAsText(file, "UTF-8");
  });
}

/**
 * 生成示例 CSV 模板
 */
export function generateUserCSVTemplate(): string {
  const headers = ["username", "email", "password", "full_name", "status"];
  const examples = [
    ["user1", "user1@example.com", "Password123", "张三", "active"],
    ["user2", "user2@example.com", "Password456", "李四", "active"],
  ];

  const bom = "\uFEFF";
  const content = [headers.join(","), ...examples.map((row) => row.join(","))].join("\n");
  return bom + content;
}
