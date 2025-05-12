//
//  mobileUITests.swift
//  mobileUITests
//
//  Created by Богдан Тарченко on 12.05.2025.
//

import XCTest

extension XCUIElement {
    func clearAndEnterText(_ text: String) {
        guard let stringValue = self.value as? String else {
            XCTFail("Tried to clear and enter text into non-string value")
            return
        }
        
        self.tap()
        
        let deleteString = String(repeating: XCUIKeyboardKey.delete.rawValue, count: stringValue.count)
        self.typeText(deleteString)
        self.typeText(text)
    }
}

final class mobileUITests: XCTestCase {
    
    override func setUpWithError() throws {
        // Put setup code here. This method is called before the invocation of each test method in the class.
        
        // In UI tests it is usually best to stop immediately when a failure occurs.
        continueAfterFailure = false
        
        // In UI tests it's important to set the initial state - such as interface orientation - required for your tests before they run. The setUp method is a good place to do this.
    }
    
    override func tearDownWithError() throws {
        // Put teardown code here. This method is called after the invocation of each test method in the class.
    }
    
    // Создать задачу с 3 символами в названии невозможно
    @MainActor
    func testTaskTitleValidation() throws {
        let app = XCUIApplication()
        app.launch()
        
        app.navigationBars.buttons["plus"].tap()
        
        let titleTextField = app.textFields["Название"]
        titleTextField.tap()
        titleTextField.typeText("123")
        
        XCTAssertFalse(app.staticTexts["Название должно содержать минимум 4 символа"].exists)
        
        app.navigationBars.buttons["Создать"].tap()
        
        XCTAssertTrue(app.staticTexts["Название должно содержать минимум 4 символа"].exists)
        
        titleTextField.tap()
        titleTextField.typeText("4")
        
        XCTAssertFalse(app.staticTexts["Название должно содержать минимум 4 символа"].exists)
    }
    
    // Создать задачу с 4 символами в названии, остальные поля формы оставить по умолчанию, проверить статус и приоритет
    @MainActor
    func testCreateMinimalTask() throws {
        let app = XCUIApplication()
        app.launch()
        
        app.navigationBars.buttons["plus"].tap()
        
        let titleTextField = app.textFields["Название"]
        titleTextField.tap()
        titleTextField.typeText("1234")
        
        app.navigationBars.buttons["Создать"].tap()
        
        XCTAssertTrue(app.navigationBars["Задачи"].exists)
        
        let taskRow = app.cells.firstMatch
        XCTAssertTrue(taskRow.exists)
        
        XCTAssertTrue(taskRow.staticTexts["Статус: ACTIVE"].exists)
        
        XCTAssertTrue(taskRow.staticTexts["MEDIUM"].exists)
    }
    
    // Отредактировать задачу, изменить все атрибуты, дедлайн установить на завтра, проверить что все данные изменились
    @MainActor
    func testEditTask() throws {
        let app = XCUIApplication()
        app.launch()
        
        app.navigationBars.buttons["plus"].tap()
        let titleTextField = app.textFields["Название"]
        titleTextField.tap()
        titleTextField.typeText("1234")
        app.navigationBars.buttons["Создать"].tap()
        
        let taskRow = app.cells.firstMatch
        taskRow.tap()
        
        app.buttons["Редактировать задачу"].tap()
        
        let editTitleTextField = app.textFields["Название"]
        editTitleTextField.tap()
        editTitleTextField.clearAndEnterText("Новое название")
        
        let descriptionField = app.textViews.firstMatch
        descriptionField.tap()
        descriptionField.typeText("Новое описание")
        
        app.buttons["Приоритет, Средний"].tap()
        app.buttons["Высокий"].tap()
        
        app.navigationBars.buttons["Сохранить"].tap()
        
        XCTAssertTrue(app.staticTexts["Новое название"].exists)
        XCTAssertTrue(app.staticTexts["Новое описание"].exists)
        XCTAssertTrue(app.staticTexts["HIGH"].exists)
        XCTAssertTrue(app.staticTexts["Дедлайн"].exists)
        
        app.navigationBars.buttons.element(boundBy: 0).tap()
        
        let updatedTaskRow = app.cells.firstMatch
        XCTAssertTrue(updatedTaskRow.staticTexts["Новое название"].exists)
        XCTAssertTrue(updatedTaskRow.staticTexts["Статус: ACTIVE"].exists)
        XCTAssertTrue(updatedTaskRow.staticTexts["HIGH"].exists)
    }
    
    // Отредактировать задачу, установить дедлайн на 1 минуту вперед, подождать 61 секунду, проверить что статус стал OVERDUE
    @MainActor
    func testTaskOverdue() throws {
        let app = XCUIApplication()
        app.launch()
        
        app.navigationBars.buttons["plus"].tap()
        let titleTextField = app.textFields["Название"]
        titleTextField.tap()
        titleTextField.typeText("1234")
        app.navigationBars.buttons["Создать"].tap()
        
        let taskRow = app.cells.firstMatch
        taskRow.tap()
        
        app.buttons["Редактировать задачу"].tap()
        
        let editTitleTextField = app.textFields["Название"]
        editTitleTextField.tap()
        editTitleTextField.clearAndEnterText("Задача с дедлайном")
        
        app.navigationBars.buttons["Сохранить"].tap()
        
        app.buttons["Редактировать задачу"].tap()

        let timeButton = app.buttons.matching(NSPredicate(format: "label CONTAINS ':'")).firstMatch
        timeButton.tap()
        
        let calendar = Calendar.current
        let now = Date()
        let oneMinuteLater = calendar.date(byAdding: .minute, value: 1, to: now)!
        
        let timeFormatter = DateFormatter()
        timeFormatter.dateFormat = "HH:mm"
        timeFormatter.locale = Locale(identifier: "en_US_POSIX")
        let timeString = timeFormatter.string(from: oneMinuteLater)
        let components = timeString.split(separator: ":")
        let minutes = String(components[1])
        
        let minuteWheel = app.pickerWheels.element(boundBy: 1)
        minuteWheel.adjust(toPickerWheelValue: minutes)
        
        app.tap()
        
        let dateButton = app.buttons.matching(NSPredicate(format: "label CONTAINS 'May'")).firstMatch
        dateButton.tap()
        
        let may12Button = app.buttons.matching(NSPredicate(format: "label CONTAINS '12 May'")).firstMatch
        may12Button.tap()
        
        let saveButton = app.navigationBars.buttons["Сохранить"]
        while !saveButton.isHittable {
            app.swipeUp()
        }
        
        saveButton.tap()
        
        app.navigationBars.buttons.element(boundBy: 0).tap()
        
        Thread.sleep(forTimeInterval: 100)
        
        app.navigationBars.buttons.element(boundBy: 0).tap()
        
        app.buttons["Применить"].tap()
        
        let updatedTaskRow = app.cells.firstMatch
        XCTAssertTrue(updatedTaskRow.staticTexts["Статус: OVERDUE"].exists)
        
        updatedTaskRow.tap()
        
        app.buttons["Отметить как выполненную"].tap()
        
        XCTAssertTrue(app.staticTexts["LATE"].exists)
        
        app.buttons["Отметить как невыполненную"].tap()
        
        XCTAssertTrue(app.staticTexts["OVERDUE"].exists)
        
        app.buttons["Отметить как выполненную"].tap()
        
        XCTAssertTrue(app.staticTexts["LATE"].exists)
    }
    
    // Создать задачу с макросом !1 в названии, проверить приоритет
    @MainActor
    func testTaskWithPriorityMacro() throws {
        let app = XCUIApplication()
        app.launch()
        
        app.navigationBars.buttons["plus"].tap()
        
        let titleTextField = app.textFields["Название"]
        titleTextField.tap()
        titleTextField.typeText("Задача !1")
        
        app.navigationBars.buttons["Создать"].tap()
        
        let taskRow = app.cells.firstMatch
        XCTAssertTrue(taskRow.staticTexts["CRITICAL"].exists)
    }
    
    // Создать задачу с макросом !before в названии, проверить дедлайн
    @MainActor
    func testTaskWithDeadlineMacro() throws {
        let app = XCUIApplication()
        app.launch()
        
        app.navigationBars.buttons["plus"].tap()
        
        let titleTextField = app.textFields["Название"]
        titleTextField.tap()
        titleTextField.typeText("Задача !before 24.04.2026")
        
        app.navigationBars.buttons["Создать"].tap()
        
        let taskRow = app.cells.firstMatch
        XCTAssertTrue(taskRow.staticTexts["Дедлайн: 24 апреля 2026 07:00"].exists)
    }
    
    // Создать задачу с обоими макросами в названии, проверить приоритет и дедлайн
    @MainActor
    func testTaskWithBothMacros() throws {
        let app = XCUIApplication()
        app.launch()
        
        app.navigationBars.buttons["plus"].tap()
        
        let titleTextField = app.textFields["Название"]
        titleTextField.tap()
        titleTextField.typeText("Задача !1 !before 24.04.2026")
        
        app.navigationBars.buttons["Создать"].tap()
        
        let taskRow = app.cells.firstMatch
        XCTAssertTrue(taskRow.staticTexts["CRITICAL"].exists)
        XCTAssertTrue(taskRow.staticTexts["Дедлайн: 24 апреля 2026 07:00"].exists)
    }
    
    // Создать задачу с обоими макросами в названии, но выбрать другие значения в форме
    @MainActor
    func testTaskWithBothMacrosAndManualValues() throws {
        let app = XCUIApplication()
        app.launch()
        
        app.navigationBars.buttons["plus"].tap()
        
        let titleTextField = app.textFields["Название"]
        titleTextField.tap()
        titleTextField.typeText("Задача !4 !before 15.05.2050")
        
        app.navigationBars.buttons["Создать"].tap()
        
        let taskRow = app.cells.firstMatch
        XCTAssertTrue(taskRow.staticTexts["LOW"].exists)
        XCTAssertTrue(taskRow.staticTexts["Дедлайн: 15 мая 2050 07:00"].exists)
        
        taskRow.tap()
        
        app.buttons["Редактировать задачу"].tap()
        
        let editTitleTextField = app.textFields["Название"]
        editTitleTextField.tap()
        editTitleTextField.clearAndEnterText("Задача !1 !before 25.04.2044")
        
        app.navigationBars.buttons["Сохранить"].tap()
        app.navigationBars.buttons.element(boundBy: 0).tap()
        
        XCTAssertTrue(app.staticTexts["LOW"].exists)
        XCTAssertTrue(app.staticTexts["Дедлайн: 15 мая 2050 07:00"].exists)
    }
    
    // Удалить задачу и проверить её отсутствие
    @MainActor
    func testTaskDeletion() throws {
        let app = XCUIApplication()
        app.launch()
        
        app.navigationBars.buttons["plus"].tap()
        let titleTextField = app.textFields["Название"]
        titleTextField.tap()
        titleTextField.typeText("Задача для удаления")
        app.navigationBars.buttons["Создать"].tap()
        
        let taskRow = app.cells.firstMatch
        taskRow.swipeLeft()
        app.buttons["trash"].tap()
        
        XCTAssertFalse(app.staticTexts["Задача для удаления"].exists)
    }
}
