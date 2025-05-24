const pool = require('../config/db');

class Blog {
  static async create({ title, content, userId }) {
    const [result] = await pool.query(
      'INSERT INTO blogs (title, content, user_id) VALUES (?, ?, ?)',
      [title, content, userId]
    );
    return result.insertId;
  }

  static async findAll() {
    const [rows] = await pool.query(`
      SELECT b.*, u.username as author 
      FROM blogs b
      JOIN users u ON b.user_id = u.id
      ORDER BY b.created_at DESC
    `);
    return rows;
  }

  static async findById(id) {
    const [rows] = await pool.query(`
      SELECT b.*, u.username as author 
      FROM blogs b
      JOIN users u ON b.user_id = u.id
      WHERE b.id = ?
    `, [id]);
    return rows[0];
  }

  static async findByUserId(userId) {
    const [rows] = await pool.query('SELECT * FROM blogs WHERE user_id = ? ORDER BY created_at DESC', [userId]);
    return rows;
  }

  static async update(id, { title, content }) {
    await pool.query(
      'UPDATE blogs SET title = ?, content = ? WHERE id = ?',
      [title, content, id]
    );
  }

  static async delete(id) {
    await pool.query('DELETE FROM blogs WHERE id = ?', [id]);
  }
}

module.exports = Blog;