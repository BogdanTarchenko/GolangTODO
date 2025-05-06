import Foundation

enum SortOrder: String, CaseIterable {
    case asc = "asc"
    case desc = "desc"
    
    var localizedName: String {
        switch self {
        case .asc:
            return "По возрастанию"
        case .desc:
            return "По убыванию"
        }
    }
}
