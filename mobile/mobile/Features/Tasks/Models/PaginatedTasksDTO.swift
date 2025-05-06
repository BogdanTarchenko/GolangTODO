struct PaginatedTasksDTO: Codable {
    var items: [Task]
    var meta: PaginationMeta
}

struct PaginationMeta: Codable {
    var page: Int
    var pageSize: Int
    var total: Int
    var totalPages: Int
    
    enum CodingKeys: String, CodingKey {
        case page
        case pageSize = "page_size"
        case total
        case totalPages = "total_pages"
    }
}
