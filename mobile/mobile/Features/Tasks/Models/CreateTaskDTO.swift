import Foundation

struct CreateTaskDTO: Codable {
    var title: String
    var description: String?
    var deadline: String?
    var priority: TaskPriority?
}
