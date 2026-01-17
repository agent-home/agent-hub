'use client'

import { useState } from 'react'
import Link from 'next/link'
import { Search, Menu, X, User, Plus } from 'lucide-react'

export function Header() {
  const [isMenuOpen, setIsMenuOpen] = useState(false)
  const [isLoggedIn, setIsLoggedIn] = useState(false)

  return (
    <header className="bg-white border-b border-gray-200 sticky top-0 z-50">
      <div className="container mx-auto px-4">
        <div className="flex items-center justify-between h-16">
          {/* Logo */}
          <Link href="/" className="flex items-center gap-2">
            <div className="w-8 h-8 bg-gradient-to-br from-primary-500 to-accent-500 rounded-lg flex items-center justify-center">
              <span className="text-white font-bold text-lg">A</span>
            </div>
            <span className="font-bold text-xl text-gray-900">AgentHub</span>
          </Link>

          {/* 桌面导航 */}
          <nav className="hidden md:flex items-center gap-6">
            <Link href="/agents" className="text-gray-600 hover:text-gray-900 transition-colors">
              智能体
            </Link>
            <Link href="/categories" className="text-gray-600 hover:text-gray-900 transition-colors">
              分类
            </Link>
            <Link href="/docs" className="text-gray-600 hover:text-gray-900 transition-colors">
              文档
            </Link>
            <Link href="/pricing" className="text-gray-600 hover:text-gray-900 transition-colors">
              定价
            </Link>
          </nav>

          {/* 搜索框 */}
          <div className="hidden md:flex flex-1 max-w-md mx-8">
            <div className="relative w-full">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400 w-4 h-4" />
              <input
                type="text"
                placeholder="搜索智能体..."
                className="w-full pl-10 pr-4 py-2 border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent"
              />
              <kbd className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 text-xs bg-gray-100 px-1.5 py-0.5 rounded">
                ⌘K
              </kbd>
            </div>
          </div>

          {/* 用户操作 */}
          <div className="hidden md:flex items-center gap-4">
            {isLoggedIn ? (
              <>
                <Link
                  href="/agents/new"
                  className="flex items-center gap-1 text-gray-600 hover:text-gray-900"
                >
                  <Plus className="w-4 h-4" />
                  <span>发布</span>
                </Link>
                <Link href="/profile" className="flex items-center gap-2">
                  <div className="w-8 h-8 bg-gray-200 rounded-full flex items-center justify-center">
                    <User className="w-4 h-4 text-gray-600" />
                  </div>
                </Link>
              </>
            ) : (
              <>
                <Link
                  href="/login"
                  className="text-gray-600 hover:text-gray-900 transition-colors"
                >
                  登录
                </Link>
                <Link
                  href="/signup"
                  className="bg-primary-600 text-white px-4 py-2 rounded-lg hover:bg-primary-700 transition-colors"
                >
                  注册
                </Link>
              </>
            )}
          </div>

          {/* 移动端菜单按钮 */}
          <button
            className="md:hidden p-2"
            onClick={() => setIsMenuOpen(!isMenuOpen)}
          >
            {isMenuOpen ? (
              <X className="w-6 h-6 text-gray-600" />
            ) : (
              <Menu className="w-6 h-6 text-gray-600" />
            )}
          </button>
        </div>
      </div>

      {/* 移动端菜单 */}
      {isMenuOpen && (
        <div className="md:hidden bg-white border-t border-gray-200">
          <div className="container mx-auto px-4 py-4">
            <div className="mb-4">
              <div className="relative">
                <Search className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400 w-4 h-4" />
                <input
                  type="text"
                  placeholder="搜索智能体..."
                  className="w-full pl-10 pr-4 py-2 border border-gray-200 rounded-lg"
                />
              </div>
            </div>
            <nav className="flex flex-col gap-2">
              <Link href="/agents" className="py-2 text-gray-600">智能体</Link>
              <Link href="/categories" className="py-2 text-gray-600">分类</Link>
              <Link href="/docs" className="py-2 text-gray-600">文档</Link>
              <Link href="/pricing" className="py-2 text-gray-600">定价</Link>
              <hr className="my-2" />
              <Link href="/login" className="py-2 text-gray-600">登录</Link>
              <Link href="/signup" className="py-2 text-primary-600">注册</Link>
            </nav>
          </div>
        </div>
      )}
    </header>
  )
}
