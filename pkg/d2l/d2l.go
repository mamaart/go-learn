package d2l

import (
	"fmt"
	"io"
	"net/http"

	"github.com/mamaart/go-learn/internal/auth"
)

type D2L struct {
	cli *http.Client
}

func New(username, password string) (*D2L, error) {
	cli, err := auth.LoginToLearn(username, password)
	if err != nil {
		return nil, fmt.Errorf("failed to login: %s", err)
	}
	return &D2L{
		cli: cli,
	}, nil
}

func (d2l *D2L) get(path string) ([]byte, error) {
	url := fmt.Sprintf("https://learn.inside.dtu.dk%s", path)
	resp, err := d2l.cli.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to do whoami: %s", err)
	}

	x, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read data: %s", err)
	}
	return x, nil
}

func (d2l *D2L) Whoami() ([]byte, error) {
	return d2l.get("/d2l/api/lp/1.2/users/whoami")
}

func (d2l *D2L) GetEnrollments() ([]byte, error) {
	return d2l.get("/d2l/api/lp/1.2/enrollments/myenrollments")
}

func (d2l *D2L) GetNews(courseId int) ([]byte, error) {
	return d2l.get(fmt.Sprintf("/d2l/api/lp/1.4/%d/news", courseId))
}

func (d2l *D2L) Version() ([]byte, error) {
	return d2l.get("/d2l/api/versions")
}

func (d2l *D2L) GetSchema() ([]byte, error) {
	return d2l.get("/d2l/api/lp/1.2/courses/schema")
}

func (d2l *D2L) GetOrg() ([]byte, error) {
	return d2l.get("/d2l/api/lp/1.2/organization/info")
}

func (d2l *D2L) GetForums(courseId int) ([]byte, error) {
	return d2l.get(fmt.Sprintf("/d2l/api/le/1.2/%d/discussions/forums/", courseId))
}

//
//    def get_raw_courses(self): return [{
//        'id': x["OrgUnit"]['Id'],
//        'title': x["OrgUnit"]['Name']}
//        for x in self.get('/d2l/api/lp/1.2/enrollments/myenrollments/')["Items"]
//        if x["OrgUnit"]["Type"]["Id"] == 3
//    ]
//
//    def transform(self, x):
//        tokens = x['title'].split()
//        if re.search(r'spring \d{2}', x['title'], re.IGNORECASE):
//            return {
//                'title'     : ' '.join(tokens[1:-2]),
//                'learn_id'  : x['id'],
//                'course_id' : tokens[0],
//                'semester'  : 'spring',
//                'year'      : int(x['title'][-2:]) + 2000,
//            }
//        elif re.search(r'f\d{2}', x['title'], re.IGNORECASE):
//            return {
//                'title'     : ' '.join(tokens[1:-1]),
//                'learn_id'  : x['id'],
//                'course_id' : tokens[0],
//                'semester'  : 'spring',
//                'year'      : int(x['title'][-2:]) + 2000,
//            }
//        elif re.search(r'fall \d{2}', x['title'], re.IGNORECASE):
//            return {
//                'title'     : ' '.join(tokens[1:-2]),
//                'learn_id'  : x['id'],
//                'course_id' : tokens[0],
//                'semester'  : 'fall',
//                'year'      : int(x['title'][-2:]) + 2000,
//            }
//        elif re.search(r'e\d{2}', x['title'], re.IGNORECASE):
//            return {
//                'title'     : ' '.join(tokens[1:-1]),
//                'learn_id'  : x['id'],
//                'course_id' : tokens[0],
//                'semester'  : 'fall',
//                'year'      : int(x['title'][-2:]) + 2000,
//            }
//
//    def get_courses(self): return [y for x in self.get_raw_courses() if (y:=self.transform(x))]
//
//
//    def get_updates(self, course_id):
//        return self.get(f'/d2l/api/le/1.3/{course_id}/updates/myUpdates')
//
//    def get_grades(self, course_id):
//        # D2L does not use grades at the moment
//        return self.get(f'/d2l/api/le/1.2/{course_id}/grades/')
//
//
//    def get_content_root(self, course_id):
//        return self.get(f'/d2l/api/le/1.3/{course_id}/content/root/')
//
//    def get_users(self, course_id):
//        return self.get(f'/d2l/api/lp/1.2/enrollments/orgUnits/{course_id}/users/'),
//
//    def get_classlist(self, course_id):
//        return self.get(f'/d2l/api/le/1.2/{course_id}/classlist/'),
//
//    def get_sections(course_id):
//        return self.get(f'/d2l/api/lp/1.2/{course_id}/sections/'),
//
//    def get_group_categories(self, course_id):
//        return self.get(f'/d2l/api/lp/1.2/{course_id}/groupcategories/')
//
//    def get_groups(self, course_no, group_category_id):
//        return self.get(f'/d2l/api/lp/1.2/{course_no}/groupcategories/{group_category_id}/groups/')
//
//    def get_group(self, course_no, group_category_id, group_id):
//        return self.get(f'/d2l/api/lp/1.2/{course_no}/groupcategories/{group_category_id}/groups/{group_id}')
//
//    def get_all_groups(self, course_no):
//        groups = self.get_group_categories(course_no)
//        group_id = groups[0]['GroupCategoryId']
//        return self.get_groups(course_no, group_id)
//
//    def add_student_to_group(self, course_no, group_no, group_id, student_id): return self.post(
//        path = f'/d2l/api/lp/1.2/{course_no}/groupcategories/{group_no}/groups/{group_id}/enrollments/',
//        data = {
//            'UserId':student_id,
//        },
//    )
//
//    def add_group(self, course_no, group_no, name, code): return self.post(
//        path = f'/d2l/api/lp/1.2/{course_no}/groupcategories/{group_no}/groups/',
//        data = {
//            'Name'          : name,
//            'Code'          : code,
//            'Description'   : {
//                'Text':'',
//                'Html':'',
//            },
//        },
//    )
