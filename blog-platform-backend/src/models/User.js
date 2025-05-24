const pool = require('../config/db');

class User {
  static async create({ username, email, password }) {
    const [result] = await pool.query(
      'INSERT INTO users (username, email, password) VALUES (?, ?, ?)',
      [username, email, password]
    );
    return result.insertId;
  }

  static async findByEmail(email) {
    const [rows] = await pool.query('SELECT * FROM users WHERE email = ?', [email]);
    return rows[0];
  }

  static async findById(id) {
    const [rows] = await pool.query('SELECT id, username, email, created_at FROM users WHERE id = ?', [id]);
    return rows[0];
  }
}

module.exports = User;