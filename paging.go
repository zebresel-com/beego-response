package beego_response

/**
 * @apiDefine Pagination
 *
 * @apiParam {Number} [skip=0] The offset to start with searching
 * @apiParam {Number} [limit=10] Maximum records for one page
 *
 *
 * @apiSuccess {Object} pagination The pagination object
 * @apiSuccess {String} pagination.first Path to first page including other params
 * @apiSuccess {String} pagination.last Path to last page including other params
 * @apiSuccess {String} pagination.previous Path to previous page
 * @apiSuccess {String} pagination.next Path to next page
 * @apiSuccess {Number} pagination.totalRecords All records on all pages
 * @apiSuccess {Number} pagination.limit Records on the current page
 * @apiSuccess {Number} pagination.pages Amount of all pages
 * @apiSuccess {Number} pagination.currentPage Number of the current page
 *
 */

type Paging struct {
	First        string `json:"first"`
	Last         string `json:"last"`
	Previous     string `json:"previous"`
	Next         string `json:"next"`
	RecordsTotal int    `json:"totalRecords"`
	RecordsPage  int    `json:"limit"`
	Pages        int    `json:"pages"`
	CurrentPage  int    `json:"currentPage"`
}
