import Link from 'next/link'
import { Download, Heart } from 'lucide-react'

interface Agent {
  id: string
  name: string
  namespace: string
  description: string
  category: string
  tags: string[]
  downloads: number
  likes: number
}

export function AgentCard({ agent }: { agent: Agent }) {
  const categoryColors: Record<string, string> = {
    coding: 'bg-blue-100 text-blue-700',
    writing: 'bg-green-100 text-green-700',
    research: 'bg-purple-100 text-purple-700',
    assistant: 'bg-orange-100 text-orange-700',
    analysis: 'bg-cyan-100 text-cyan-700',
    creative: 'bg-pink-100 text-pink-700',
    education: 'bg-yellow-100 text-yellow-700',
    business: 'bg-indigo-100 text-indigo-700',
    tooling: 'bg-gray-100 text-gray-700',
    other: 'bg-gray-100 text-gray-700',
  }

  const formatNumber = (num: number) => {
    if (num >= 1000) {
      return (num / 1000).toFixed(1) + 'k'
    }
    return num.toString()
  }

  return (
    <Link href={`/${agent.namespace}/${agent.name}`}>
      <div className="bg-white rounded-xl border border-gray-200 p-5 card-hover cursor-pointer">
        {/* 头部 */}
        <div className="flex items-start justify-between mb-3">
          <div className="flex items-center gap-2">
            <div className="w-10 h-10 bg-gradient-to-br from-primary-400 to-accent-400 rounded-lg flex items-center justify-center text-white font-bold">
              {agent.name.charAt(0).toUpperCase()}
            </div>
            <div>
              <div className="font-semibold text-gray-900">{agent.name}</div>
              <div className="text-sm text-gray-500">{agent.namespace}</div>
            </div>
          </div>
        </div>

        {/* 描述 */}
        <p className="text-gray-600 text-sm mb-4 line-clamp-2">
          {agent.description}
        </p>

        {/* 标签 */}
        <div className="flex flex-wrap gap-1 mb-4">
          <span className={`text-xs px-2 py-0.5 rounded-full ${categoryColors[agent.category] || categoryColors.other}`}>
            {agent.category}
          </span>
          {agent.tags.slice(0, 2).map((tag) => (
            <span key={tag} className="text-xs px-2 py-0.5 rounded-full bg-gray-100 text-gray-600">
              {tag}
            </span>
          ))}
        </div>

        {/* 统计 */}
        <div className="flex items-center gap-4 text-sm text-gray-500">
          <div className="flex items-center gap-1">
            <Download className="w-4 h-4" />
            <span>{formatNumber(agent.downloads)}</span>
          </div>
          <div className="flex items-center gap-1">
            <Heart className="w-4 h-4" />
            <span>{formatNumber(agent.likes)}</span>
          </div>
        </div>
      </div>
    </Link>
  )
}
