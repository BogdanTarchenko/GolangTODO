struct UpdateTaskStatusDTO: Codable {
    let isCompleted: Bool
    
    enum CodingKeys: String, CodingKey {
        case isCompleted = "is_completed"
    }
}
