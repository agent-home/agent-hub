'use client'

import { useState, useEffect } from 'react'
import Link from 'next/link'
import { Search, ArrowRight, Download, Heart, Zap, Code, BookOpen, Sparkles } from 'lucide-react'
import { AgentCard } from '@/components/AgentCard'
import { CategoryCard } from '@/components/CategoryCard'

// 模拟数据
const featuredAgents = [
  {
    id: '1',
    name: 'code-reviewer',
    namespace: 'agenthub',
    description: '专业的代码审查智能体，支持多种编程语言，提供代码质量分析和安全漏洞检测',
    category: 'coding',
    tags: ['coding', 'review', 'security'],
    downloads: 12580,
    likes: 892,
  },
  {
    id: '2',
    name: 'research-assistant',
    namespace: 'agenthub',
    description: '学术研究助手，帮助你搜索文献、总结论文、生成研究报告',
    category: 'research',
    tags: ['research', 'academic', 'paper'],
    downloads: 8920,
    likes: 654,
  },
  {
    id: '3',
    name: 'sql-expert',
    namespace: 'datamaster',
    description: '自然语言转 SQL 专家，支持多种数据库，生成优化的查询语句',
    category: 'coding',
    tags: ['sql', 'database', 'query'],
    downloads: 7650,
    likes: 521,
  },
  {
    id: '4',
    name: 'writing-coach',
    namespace: 'creator',
    description: '写作教练，帮助改进文章结构、润色语言、提供写作建议',
    category: 'writing',
    tags: ['writing', 'editing', 'creative'],
    downloads: 6890,
    likes: 489,
  },
]

const categories = [
  { id: 'coding', name: '编程开发', icon: Code, count: 1250, color: 'bg-blue-500' },
  { id: 'writing', name: '写作创作', icon: BookOpen, count: 890, color: 'bg-green-500' },
  { id: 'research', name: '研究探索', icon: Sparkles, count: 650, color: 'bg-purple-500' },
  { id: 'assistant', name: '通用助手', icon: Zap, count: 2100, color: 'bg-orange-500' },
]

export default function Home() {
  const [searchQuery, setSearchQuery] = useState('')

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault()
    if (searchQuery.trim()) {
      window.location.href = `/search?q=${encodeURIComponent(searchQuery)}`
    }
  }

  return (
    <div>
      {/* Hero Section */}
      <section className="gradient-bg text-white py-20">
        <div className="container mx-auto px-4 text-center">
          <h1 className="text-5xl font-bold mb-6">
            智能体的开源社区
          </h1>
          <p className="text-xl mb-8 text-white/90 max-w-2xl mx-auto">
            发现、分享和运行 AI 智能体。AgentHub 是智能体领域的 Hugging Face。
          </p>
          
          {/* 搜索框 */}
          <form onSubmit={handleSearch} className="max-w-2xl mx-auto mb-8">
            <div className="relative">
              <Search className="absolute left-4 top-1/2 -translate-y-1/2 text-gray-400 w-5 h-5" />
              <input
                type="text"
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                placeholder="搜索智能体... 例如: code review, 翻译助手"
                className="w-full pl-12 pr-32 py-4 rounded-full text-gray-900 text-lg focus:outline-none focus:ring-4 focus:ring-white/30"
              />
              <button
                type="submit"
                className="absolute right-2 top-1/2 -translate-y-1/2 bg-primary-600 hover:bg-primary-700 text-white px-6 py-2 rounded-full transition-colors"
              >
                搜索
              </button>
            </div>
          </form>

          {/* 统计 */}
          <div className="flex justify-center gap-12 text-white/90">
            <div>
              <div className="text-3xl font-bold">10,000+</div>
              <div className="text-sm">智能体</div>
            </div>
            <div>
              <div className="text-3xl font-bold">50,000+</div>
              <div className="text-sm">开发者</div>
            </div>
            <div>
              <div className="text-3xl font-bold">1M+</div>
              <div className="text-sm">调用次数</div>
            </div>
          </div>
        </div>
      </section>

      {/* 分类 */}
      <section className="py-12 bg-white">
        <div className="container mx-auto px-4">
          <div className="flex items-center justify-between mb-8">
            <h2 className="text-2xl font-bold text-gray-900">浏览分类</h2>
            <Link href="/categories" className="text-primary-600 hover:text-primary-700 flex items-center gap-1">
              查看全部 <ArrowRight className="w-4 h-4" />
            </Link>
          </div>
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
            {categories.map((category) => (
              <CategoryCard key={category.id} category={category} />
            ))}
          </div>
        </div>
      </section>

      {/* 热门智能体 */}
      <section className="py-12">
        <div className="container mx-auto px-4">
          <div className="flex items-center justify-between mb-8">
            <h2 className="text-2xl font-bold text-gray-900">热门智能体</h2>
            <Link href="/agents" className="text-primary-600 hover:text-primary-700 flex items-center gap-1">
              查看全部 <ArrowRight className="w-4 h-4" />
            </Link>
          </div>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
            {featuredAgents.map((agent) => (
              <AgentCard key={agent.id} agent={agent} />
            ))}
          </div>
        </div>
      </section>

      {/* CLI 安装 */}
      <section className="py-12 bg-gray-900 text-white">
        <div className="container mx-auto px-4">
          <div className="max-w-3xl mx-auto text-center">
            <h2 className="text-3xl font-bold mb-4">用命令行管理智能体</h2>
            <p className="text-gray-400 mb-8">
              安装 AgentHub CLI，一行命令搜索、下载、运行智能体
            </p>
            <div className="bg-gray-800 rounded-lg p-6 text-left">
              <div className="text-gray-400 text-sm mb-2"># 安装 CLI</div>
              <div className="text-green-400 font-mono mb-4">
                $ go install github.com/agenthub/cli@latest
              </div>
              <div className="text-gray-400 text-sm mb-2"># 搜索智能体</div>
              <div className="text-green-400 font-mono mb-4">
                $ agenthub search "code review"
              </div>
              <div className="text-gray-400 text-sm mb-2"># 运行智能体</div>
              <div className="text-green-400 font-mono">
                $ agenthub run agenthub/code-reviewer
              </div>
            </div>
            <div className="mt-8 flex justify-center gap-4">
              <Link
                href="/docs/cli"
                className="bg-white text-gray-900 px-6 py-3 rounded-lg font-medium hover:bg-gray-100 transition-colors"
              >
                阅读文档
              </Link>
              <Link
                href="/docs/quickstart"
                className="border border-white/30 text-white px-6 py-3 rounded-lg font-medium hover:bg-white/10 transition-colors"
              >
                快速开始
              </Link>
            </div>
          </div>
        </div>
      </section>

      {/* CTA */}
      <section className="py-16 bg-white">
        <div className="container mx-auto px-4 text-center">
          <h2 className="text-3xl font-bold mb-4">准备好发布你的智能体了吗？</h2>
          <p className="text-gray-600 mb-8 max-w-xl mx-auto">
            加入 AgentHub 社区，分享你的智能体，让更多人使用
          </p>
          <div className="flex justify-center gap-4">
            <Link
              href="/signup"
              className="bg-primary-600 text-white px-8 py-3 rounded-lg font-medium hover:bg-primary-700 transition-colors"
            >
              免费注册
            </Link>
            <Link
              href="/docs/publish"
              className="border border-gray-300 text-gray-700 px-8 py-3 rounded-lg font-medium hover:bg-gray-50 transition-colors"
            >
              发布指南
            </Link>
          </div>
        </div>
      </section>
    </div>
  )
}
