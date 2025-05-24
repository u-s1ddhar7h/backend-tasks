const pool = require('../config/db');

class Comment {
  static async create({ content, userId, blogId }) {
    const [result] = await pool.query(
      'INSERT INTO comments (content, user_id, blog_id) VALUES (?, ?, ?)',
      [content, userId, blogId]
    );
    return result.insertId;
  }

  static async findByBlogId(blogId) {
    const [rows] = await pool.query(`
      SELECT c.*, u.username as author 
      FROM comments c
      JOIN users u ON c.user_id = u.id
      WHERE c.blog_id = ?
      ORDER BY c.created_at DESC
    `, [blogId]);
    return rows;
  }

  static async delete(id) {
    await pool.query('DELETE FROM comments WHERE id = ?', [id]);
  }
}

module.exports = Comment;