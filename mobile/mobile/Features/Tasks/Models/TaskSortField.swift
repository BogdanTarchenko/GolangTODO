import Foundation

enum TaskSortField: String, CaseIterable {
    case createdAt = "created_at"
    case deadline = "deadline"
    case priority = "priority"
    
    var localizedName: String {
        switch self {
        case .createdAt:
            return "Дата создания"
        case .deadline:
            return "Дедлайн"
        case .priority:
            return "Приоритет"
        }
    }
}
