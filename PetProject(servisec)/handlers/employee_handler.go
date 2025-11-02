package handlers

import (
	"database/sql"
	"employee-service/models"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type EmployeeHandler struct {
	DB *sql.DB
}

func NewEmployeeHandler(db *sql.DB) *EmployeeHandler {
	return &EmployeeHandler{DB: db}
}

//  Добавление сотрудника
func (h *EmployeeHandler) AddEmployee(w http.ResponseWriter, r *http.Request) {
	var employee models.Employee
	if err := json.NewDecoder(r.Body).Decode(&employee); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Валидация обязательных полей
	if employee.Name == "" || employee.Surname == "" {
		http.Error(w, "Name and surname are required", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO employees (name, surname, phone, company_id, passport_type,
              passport_number, department_name, department_phone)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`

	var id int
	err := h.DB.QueryRow(query,
		employee.Name,
		employee.Surname,
		employee.Phone,
		employee.CompanyId,
		employee.Passport.Type,
		employee.Passport.Number,
		employee.Department.Name,
		employee.Department.Phone).Scan(&id)

	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": id})
}

// Удаление сотрудника по ID
func (h *EmployeeHandler) DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid employee ID", http.StatusBadRequest)
		return
	}

	result, err := h.DB.Exec("DELETE FROM employees WHERE id = $1", id)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Employee not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Employee deleted successfully"})
}

// Список сотрудников по компании
func (h *EmployeeHandler) GetEmployeesByCompany(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	companyId, err := strconv.Atoi(vars["companyId"])
	if err != nil {
		http.Error(w, "Invalid company ID", http.StatusBadRequest)
		return
	}

	rows, err := h.DB.Query(`
        SELECT id, name, surname, phone, company_id, passport_type,
               passport_number, department_name, department_phone
        FROM employees WHERE company_id = $1`, companyId)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var employees []models.Employee
	for rows.Next() {
		var emp models.Employee
		err := rows.Scan(
			&emp.Id, &emp.Name, &emp.Surname, &emp.Phone, &emp.CompanyId,
			&emp.Passport.Type, &emp.Passport.Number,
			&emp.Department.Name, &emp.Department.Phone,
		)
		if err != nil {
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		employees = append(employees, emp)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(employees)
}

// Список сотрудников по отделу
func (h *EmployeeHandler) GetEmployeesByDepartment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	companyId, err := strconv.Atoi(vars["companyId"])
	if err != nil {
		http.Error(w, "Invalid company ID", http.StatusBadRequest)
		return
	}

	departmentName := vars["departmentName"]

	rows, err := h.DB.Query(`
        SELECT id, name, surname, phone, company_id, passport_type,
               passport_number, department_name, department_phone
        FROM employees WHERE company_id = $1 AND department_name = $2`,
		companyId, departmentName)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var employees []models.Employee
	for rows.Next() {
		var emp models.Employee
		err := rows.Scan(
			&emp.Id, &emp.Name, &emp.Surname, &emp.Phone, &emp.CompanyId,
			&emp.Passport.Type, &emp.Passport.Number,
			&emp.Department.Name, &emp.Department.Phone,
		)
		if err != nil {
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		employees = append(employees, emp)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(employees)
}

// Обновление сотрудника по ID
func (h *EmployeeHandler) UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid employee ID", http.StatusBadRequest)
		return
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Проверяем существование сотрудника
	var exists bool
	err = h.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM employees WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(w, "Employee not found", http.StatusNotFound)
		return
	}

	// Динамическое построение запроса
	query := "UPDATE employees SET "
	params := []interface{}{}
	paramCount := 1

	fields := map[string]string{
		"name": "name", "surname": "surname", "phone": "phone",
		"companyId": "company_id", "passport_type": "passport_type",
		"passport_number": "passport_number", "department_name": "department_name",
		"department_phone": "department_phone",
	}

	first := true
	for jsonField, dbField := range fields {
		if value, ok := updates[jsonField]; ok {
			if !first {
				query += ", "
			}
			query += dbField + " = $" + strconv.Itoa(paramCount)
			params = append(params, value)
			paramCount++
			first = false
		}
	}

	if len(params) == 0 {
		http.Error(w, "No fields to update", http.StatusBadRequest)
		return
	}

	query += " WHERE id = $" + strconv.Itoa(paramCount)
	params = append(params, id)

	_, err = h.DB.Exec(query, params...)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Employee updated successfully"})
}
