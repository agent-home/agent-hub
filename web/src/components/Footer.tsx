import Link from 'next/link'
import { Github, Twitter } from 'lucide-react'

export function Footer() {
  return (
    <footer className="bg-gray-900 text-gray-400">
      <div className="container mx-auto px-4 py-12">
        <div className="grid grid-cols-2 md:grid-cols-4 gap-8">
          {/* 产品 */}
          <div>
            <h3 className="text-white font-semibold mb-4">产品</h3>
            <ul className="space-y-2">
              <li><Link href="/agents" className="hover:text-white transition-colors">智能体市场</Link></li>
              <li><Link href="/playground" className="hover:text-white transition-colors">Playground</Link></li>
              <li><Link href="/pricing" className="hover:text-white transition-colors">定价</Link></li>
              <li><Link href="/enterprise" className="hover:text-white transition-colors">企业版</Link></li>
            </ul>
          </div>

          {/* 开发者 */}
          <div>
            <h3 className="text-white font-semibold mb-4">开发者</h3>
            <ul className="space-y-2">
              <li><Link href="/docs" className="hover:text-white transition-colors">文档</Link></li>
              <li><Link href="/docs/api" className="hover:text-white transition-colors">API 参考</Link></li>
              <li><Link href="/docs/cli" className="hover:text-white transition-colors">CLI 工具</Link></li>
              <li><Link href="/docs/sdk" className="hover:text-white transition-colors">SDK</Link></li>
            </ul>
          </div>

          {/* 社区 */}
          <div>
            <h3 className="text-white font-semibold mb-4">社区</h3>
            <ul className="space-y-2">
              <li><Link href="/blog" className="hover:text-white transition-colors">博客</Link></li>
              <li><Link href="/discord" className="hover:text-white transition-colors">Discord</Link></li>
              <li><Link href="/github" className="hover:text-white transition-colors">GitHub</Link></li>
              <li><Link href="/changelog" className="hover:text-white transition-colors">更新日志</Link></li>
            </ul>
          </div>

          {/* 公司 */}
          <div>
            <h3 className="text-white font-semibold mb-4">关于</h3>
            <ul className="space-y-2">
              <li><Link href="/about" className="hover:text-white transition-colors">关于我们</Link></li>
              <li><Link href="/careers" className="hover:text-white transition-colors">加入我们</Link></li>
              <li><Link href="/privacy" className="hover:text-white transition-colors">隐私政策</Link></li>
              <li><Link href="/terms" className="hover:text-white transition-colors">服务条款</Link></li>
            </ul>
          </div>
        </div>

        <div className="border-t border-gray-800 mt-12 pt-8 flex flex-col md:flex-row items-center justify-between">
          <div className="flex items-center gap-2 mb-4 md:mb-0">
            <div className="w-6 h-6 bg-gradient-to-br from-primary-500 to-accent-500 rounded flex items-center justify-center">
              <span className="text-white font-bold text-sm">A</span>
            </div>
            <span className="text-white font-semibold">AgentHub</span>
          </div>

          <div className="text-sm mb-4 md:mb-0">
            © 2026 AgentHub. All rights reserved.
          </div>

          <div className="flex items-center gap-4">
            <a href="https://github.com/agenthub" target="_blank" rel="noopener noreferrer" className="hover:text-white transition-colors">
              <Github className="w-5 h-5" />
            </a>
            <a href="https://twitter.com/agenthub" target="_blank" rel="noopener noreferrer" className="hover:text-white transition-colors">
              <Twitter className="w-5 h-5" />
            </a>
          </div>
        </div>
      </div>
    </footer>
  )
}
