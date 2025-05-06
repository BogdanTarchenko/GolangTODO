import Foundation

struct Task: Identifiable, Codable {
    let id: String
    var title: String
    var description: String?
    var deadline: String?
    var status: TaskStatus
    var priority: TaskPriority
    let createdAt: String
    var updatedAt: String?
    
    enum CodingKeys: String, CodingKey {
        case id
        case title
        case description
        case deadline
        case status
        case priority
        case createdAt = "created_at"
        case updatedAt = "updated_at"
    }
}
