import Link from 'next/link'
import { LucideIcon } from 'lucide-react'

interface Category {
  id: string
  name: string
  icon: LucideIcon
  count: number
  color: string
}

export function CategoryCard({ category }: { category: Category }) {
  const Icon = category.icon

  return (
    <Link href={`/categories/${category.id}`}>
      <div className="bg-white rounded-xl border border-gray-200 p-5 card-hover cursor-pointer">
        <div className="flex items-center gap-3">
          <div className={`w-10 h-10 ${category.color} rounded-lg flex items-center justify-center`}>
            <Icon className="w-5 h-5 text-white" />
          </div>
          <div>
            <div className="font-semibold text-gray-900">{category.name}</div>
            <div className="text-sm text-gray-500">{category.count} 个智能体</div>
          </div>
        </div>
      </div>
    </Link>
  )
}
