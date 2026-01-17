'use client'

import { useState } from 'react'
import Link from 'next/link'
import { 
  Download, Heart, Share2, Code, Play, FileText, 
  GitBranch, Clock, User, ExternalLink, Copy, Check,
  MessageSquare, Send
} from 'lucide-react'

// 模拟数据
const agentData = {
  name: 'code-reviewer',
  namespace: 'agenthub',
  description: '专业的代码审查智能体，支持多种编程语言，提供代码质量分析、安全漏洞检测和优化建议',
  longDescription: `
## 功能特性

- **多语言支持**: 支持 Python、JavaScript、TypeScript、Go、Java、C++ 等主流语言
- **安全检测**: 自动检测 SQL 注入、XSS、CSRF 等安全漏洞
- **代码质量**: 分析代码复杂度、重复代码、命名规范等
- **优化建议**: 提供性能优化和重构建议

## 使用场景

1. 代码提交前的自动审查
2. Pull Request 审查助手
3. 代码库安全扫描
4. 新人代码培训
`,
  category: 'coding',
  tags: ['coding', 'review', 'security', 'quality'],
  license: 'Apache-2.0',
  downloads: 12580,
  likes: 892,
  versions: [
    { version: '1.2.0', date: '2026-01-15', isLatest: true },
    { version: '1.1.0', date: '2026-01-01' },
    { version: '1.0.0', date: '2025-12-15' },
  ],
  author: {
    name: 'agenthub',
    avatar: null,
  },
  createdAt: '2025-12-15',
  updatedAt: '2026-01-15',
}

export default function AgentPage({ params }: { params: { namespace: string; name: string } }) {
  const [activeTab, setActiveTab] = useState<'readme' | 'versions' | 'playground'>('readme')
  const [copied, setCopied] = useState(false)
  const [message, setMessage] = useState('')
  const [chatHistory, setChatHistory] = useState<{role: string; content: string}[]>([])
  const [isLoading, setIsLoading] = useState(false)

  const fullName = `${params.namespace}/${params.name}`

  const copyCommand = () => {
    navigator.clipboard.writeText(`agenthub pull ${fullName}`)
    setCopied(true)
    setTimeout(() => setCopied(false), 2000)
  }

  const sendMessage = async () => {
    if (!message.trim() || isLoading) return

    const userMessage = message
    setMessage('')
    setChatHistory(prev => [...prev, { role: 'user', content: userMessage }])
    setIsLoading(true)

    // 模拟响应
    setTimeout(() => {
      setChatHistory(prev => [...prev, { 
        role: 'assistant', 
        content: '这是一个演示响应。在实际实现中，这里会调用智能体 API 获取真实响应。' 
      }])
      setIsLoading(false)
    }, 1000)
  }

  return (
    <div className="min-h-screen bg-gray-50">
      {/* 头部 */}
      <div className="bg-white border-b border-gray-200">
        <div className="container mx-auto px-4 py-6">
          <div className="flex items-start justify-between">
            <div className="flex items-start gap-4">
              <div className="w-16 h-16 bg-gradient-to-br from-primary-400 to-accent-400 rounded-xl flex items-center justify-center text-white text-2xl font-bold">
                {agentData.name.charAt(0).toUpperCase()}
              </div>
              <div>
                <div className="flex items-center gap-2 mb-1">
                  <h1 className="text-2xl font-bold text-gray-900">{agentData.name}</h1>
                  <span className="text-xs bg-green-100 text-green-700 px-2 py-0.5 rounded-full">
                    v{agentData.versions[0].version}
                  </span>
                </div>
                <Link href={`/${params.namespace}`} className="text-gray-500 hover:text-primary-600">
                  {params.namespace}
                </Link>
                <p className="text-gray-600 mt-2 max-w-2xl">{agentData.description}</p>
                <div className="flex flex-wrap gap-2 mt-3">
                  {agentData.tags.map((tag) => (
                    <span key={tag} className="text-xs px-2 py-1 bg-gray-100 text-gray-600 rounded-full">
                      {tag}
                    </span>
                  ))}
                </div>
              </div>
            </div>

            <div className="flex items-center gap-3">
              <button className="flex items-center gap-2 px-4 py-2 border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors">
                <Heart className="w-4 h-4" />
                <span>{agentData.likes}</span>
              </button>
              <button className="flex items-center gap-2 px-4 py-2 border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors">
                <Share2 className="w-4 h-4" />
                <span>分享</span>
              </button>
            </div>
          </div>

          {/* 统计 */}
          <div className="flex items-center gap-6 mt-6 text-sm text-gray-500">
            <div className="flex items-center gap-1">
              <Download className="w-4 h-4" />
              <span>{agentData.downloads.toLocaleString()} 下载</span>
            </div>
            <div className="flex items-center gap-1">
              <GitBranch className="w-4 h-4" />
              <span>{agentData.versions.length} 个版本</span>
            </div>
            <div className="flex items-center gap-1">
              <Clock className="w-4 h-4" />
              <span>更新于 {agentData.updatedAt}</span>
            </div>
            <div className="flex items-center gap-1">
              <FileText className="w-4 h-4" />
              <span>{agentData.license}</span>
            </div>
          </div>
        </div>
      </div>

      {/* 内容区 */}
      <div className="container mx-auto px-4 py-6">
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* 主内容 */}
          <div className="lg:col-span-2">
            {/* Tab 切换 */}
            <div className="bg-white rounded-xl border border-gray-200 overflow-hidden">
              <div className="flex border-b border-gray-200">
                <button
                  onClick={() => setActiveTab('readme')}
                  className={`px-6 py-3 text-sm font-medium transition-colors ${
                    activeTab === 'readme'
                      ? 'text-primary-600 border-b-2 border-primary-600'
                      : 'text-gray-500 hover:text-gray-700'
                  }`}
                >
                  <FileText className="w-4 h-4 inline mr-2" />
                  说明
                </button>
                <button
                  onClick={() => setActiveTab('versions')}
                  className={`px-6 py-3 text-sm font-medium transition-colors ${
                    activeTab === 'versions'
                      ? 'text-primary-600 border-b-2 border-primary-600'
                      : 'text-gray-500 hover:text-gray-700'
                  }`}
                >
                  <GitBranch className="w-4 h-4 inline mr-2" />
                  版本
                </button>
                <button
                  onClick={() => setActiveTab('playground')}
                  className={`px-6 py-3 text-sm font-medium transition-colors ${
                    activeTab === 'playground'
                      ? 'text-primary-600 border-b-2 border-primary-600'
                      : 'text-gray-500 hover:text-gray-700'
                  }`}
                >
                  <Play className="w-4 h-4 inline mr-2" />
                  试用
                </button>
              </div>

              <div className="p-6">
                {activeTab === 'readme' && (
                  <div className="prose max-w-none">
                    <div dangerouslySetInnerHTML={{ __html: agentData.longDescription.replace(/\n/g, '<br/>') }} />
                  </div>
                )}

                {activeTab === 'versions' && (
                  <div className="space-y-3">
                    {agentData.versions.map((v) => (
                      <div
                        key={v.version}
                        className="flex items-center justify-between p-4 bg-gray-50 rounded-lg"
                      >
                        <div className="flex items-center gap-3">
                          <span className="font-mono font-medium">{v.version}</span>
                          {v.isLatest && (
                            <span className="text-xs bg-green-100 text-green-700 px-2 py-0.5 rounded-full">
                              最新
                            </span>
                          )}
                        </div>
                        <div className="flex items-center gap-4">
                          <span className="text-sm text-gray-500">{v.date}</span>
                          <button className="text-primary-600 hover:text-primary-700 text-sm">
                            下载
                          </button>
                        </div>
                      </div>
                    ))}
                  </div>
                )}

                {activeTab === 'playground' && (
                  <div className="h-96 flex flex-col">
                    {/* 聊天历史 */}
                    <div className="flex-1 overflow-y-auto space-y-4 mb-4">
                      {chatHistory.length === 0 ? (
                        <div className="text-center text-gray-500 py-12">
                          <MessageSquare className="w-12 h-12 mx-auto mb-4 text-gray-300" />
                          <p>输入消息开始与智能体对话</p>
                        </div>
                      ) : (
                        chatHistory.map((msg, i) => (
                          <div
                            key={i}
                            className={`flex ${msg.role === 'user' ? 'justify-end' : 'justify-start'}`}
                          >
                            <div
                              className={`max-w-[80%] px-4 py-2 rounded-lg ${
                                msg.role === 'user'
                                  ? 'bg-primary-600 text-white'
                                  : 'bg-gray-100 text-gray-900'
                              }`}
                            >
                              {msg.content}
                            </div>
                          </div>
                        ))
                      )}
                      {isLoading && (
                        <div className="flex justify-start">
                          <div className="bg-gray-100 text-gray-900 px-4 py-2 rounded-lg">
                            <span className="animate-pulse">正在思考...</span>
                          </div>
                        </div>
                      )}
                    </div>

                    {/* 输入框 */}
                    <div className="flex gap-2">
                      <input
                        type="text"
                        value={message}
                        onChange={(e) => setMessage(e.target.value)}
                        onKeyDown={(e) => e.key === 'Enter' && sendMessage()}
                        placeholder="输入消息..."
                        className="flex-1 px-4 py-2 border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500"
                      />
                      <button
                        onClick={sendMessage}
                        disabled={isLoading}
                        className="px-4 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700 disabled:opacity-50 transition-colors"
                      >
                        <Send className="w-5 h-5" />
                      </button>
                    </div>
                  </div>
                )}
              </div>
            </div>
          </div>

          {/* 侧边栏 */}
          <div className="space-y-4">
            {/* 安装命令 */}
            <div className="bg-white rounded-xl border border-gray-200 p-4">
              <h3 className="font-semibold text-gray-900 mb-3">安装</h3>
              <div className="bg-gray-900 rounded-lg p-3 flex items-center justify-between">
                <code className="text-green-400 text-sm">
                  agenthub pull {fullName}
                </code>
                <button
                  onClick={copyCommand}
                  className="text-gray-400 hover:text-white transition-colors"
                >
                  {copied ? <Check className="w-4 h-4" /> : <Copy className="w-4 h-4" />}
                </button>
              </div>
            </div>

            {/* 运行命令 */}
            <div className="bg-white rounded-xl border border-gray-200 p-4">
              <h3 className="font-semibold text-gray-900 mb-3">运行</h3>
              <div className="bg-gray-900 rounded-lg p-3">
                <code className="text-green-400 text-sm">
                  agenthub run {fullName}
                </code>
              </div>
            </div>

            {/* 作者信息 */}
            <div className="bg-white rounded-xl border border-gray-200 p-4">
              <h3 className="font-semibold text-gray-900 mb-3">作者</h3>
              <Link href={`/${agentData.author.name}`} className="flex items-center gap-3 hover:bg-gray-50 -mx-2 px-2 py-2 rounded-lg transition-colors">
                <div className="w-10 h-10 bg-gray-200 rounded-full flex items-center justify-center">
                  <User className="w-5 h-5 text-gray-500" />
                </div>
                <div>
                  <div className="font-medium text-gray-900">{agentData.author.name}</div>
                  <div className="text-sm text-gray-500">查看更多智能体</div>
                </div>
              </Link>
            </div>

            {/* 相关链接 */}
            <div className="bg-white rounded-xl border border-gray-200 p-4">
              <h3 className="font-semibold text-gray-900 mb-3">链接</h3>
              <div className="space-y-2">
                <a href="#" className="flex items-center gap-2 text-gray-600 hover:text-primary-600 transition-colors">
                  <Code className="w-4 h-4" />
                  <span>源代码</span>
                  <ExternalLink className="w-3 h-3 ml-auto" />
                </a>
                <a href="#" className="flex items-center gap-2 text-gray-600 hover:text-primary-600 transition-colors">
                  <FileText className="w-4 h-4" />
                  <span>API 文档</span>
                  <ExternalLink className="w-3 h-3 ml-auto" />
                </a>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
