const Blog = require('../models/Blog');

const createBlog = async (req, res) => {
  try {
    const { title, content } = req.body;
    const blogId = await Blog.create({ title, content, userId: req.user.id });
    const blog = await Blog.findById(blogId);
    res.status(201).json(blog);
  } catch (error) {
    res.status(500).json({ message: error.message });
  }
};

const getAllBlogs = async (req, res) => {
  try {
    const blogs = await Blog.findAll();
    res.json(blogs);
  } catch (error) {
    res.status(500).json({ message: error.message });
  }
};

const getBlogById = async (req, res) => {
  try {
    const blog = await Blog.findById(req.params.id);
    if (!blog) {
      return res.status(404).json({ message: 'Blog not found' });
    }
    res.json(blog);
  } catch (error) {
    res.status(500).json({ message: error.message });
  }
};

const getUserBlogs = async (req, res) => {
  try {
    const blogs = await Blog.findByUserId(req.user.id);
    res.json(blogs);
  } catch (error) {
    res.status(500).json({ message: error.message });
  }
};

const updateBlog = async (req, res) => {
  try {
    const { title, content } = req.body;
    const blog = await Blog.findById(req.params.id);

    if (!blog) {
      return res.status(404).json({ message: 'Blog not found' });
    }

    if (blog.user_id !== req.user.id) {
      return res.status(403).json({ message: 'Not authorized to update this blog' });
    }

    await Blog.update(req.params.id, { title, content });
    const updatedBlog = await Blog.findById(req.params.id);
    res.json(updatedBlog);
  } catch (error) {
    res.status(500).json({ message: error.message });
  }
};

const deleteBlog = async (req, res) => {
  try {
    const blog = await Blog.findById(req.params.id);

    if (!blog) {
      return res.status(404).json({ message: 'Blog not found' });
    }

    if (blog.user_id !== req.user.id) {
      return res.status(403).json({ message: 'Not authorized to delete this blog' });
    }

    await Blog.delete(req.params.id);
    res.json({ message: 'Blog deleted successfully' });
  } catch (error) {
    res.status(500).json({ message: error.message });
  }
};

module.exports = {
  createBlog,
  getAllBlogs,
  getBlogById,
  getUserBlogs,
  updateBlog,
  deleteBlog
};