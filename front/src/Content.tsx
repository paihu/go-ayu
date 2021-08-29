import axios from "axios";
import { useEffect, useState } from "react";
import { Link, useLocation } from "react-router-dom";
import { toast } from "react-toastify";
import { baseUrl } from "./const";
import "./Content.css";
import { Post } from "./Post";

const client = axios.create({ withCredentials: true, timeout: 1000 });

interface ContentProps {
  page: number;
}

type generateTopicStringProps = {
  count?: number,
  limit_count?: number,
  limit_date?: string,
  require_message?: boolean,
}

const generateTopicString: (g: generateTopicStringProps) =>  string = (g: generateTopicStringProps)=> {
  if (g.limit_count && g.count && g.limit_count !==0 && g.limit_count <= g.count){
    return "提供終了"
  }
  if (g.limit_date){
  const t = new Date();
  t.setTime(Date.parse(g.limit_date));
  if (t < new Date()){
    return "提供終了"
  }
}
  let s = "";
  if (g.limit_count && g.limit_count !==0){
    s += `あと${g.limit_count - (g.count||0)}人`
  }
  if(g.require_message){
    if (s.length !== 0 ) s += " "
    s += "要書込"
  }
  return s
}

const parseDate: (s: string) => string = (s: string) => {
  const t = new Date();
  t.setTime(Date.parse(s));
  return `${t.getFullYear()}/${(t.getMonth() + 1)
    .toString()
    .padStart(2, "0")}/${t.getDate().toString().padStart(2, "0")}`;
};
const Content: React.FC<ContentProps> = ({ page }) => {
  
  const location = useLocation();
  const [posts, setPosts] = useState([] as Post[]);
  const [delete_post_ids, setDeletePostIds] = useState([] as number[]);
  const [delete_password, setDeletePassword] = useState("");
  const deletePost: (e: React.MouseEvent<HTMLInputElement>) => void = (async e =>{
    e.preventDefault()
    console.log(delete_post_ids)
     await Promise.all(delete_post_ids.map(async (id) =>{
      const resp = await client.delete(`${baseUrl}/api/post/${id}`, {headers:{"Authorization": delete_password} }).catch(
        (e)=> {console.log(e);toast.error(`id: ${id}, delete failed`,{position: "bottom-right"})}
      )
      if (resp){
      toast.success(resp.data,{position: "bottom-right"})
      }
    }));
      const resp = await client
        .get(`${baseUrl}/api/post/page/${page}`)
        .catch((e) => {
          console.log(e);
        });
      resp ? setPosts(resp.data as Post[]) : console.log(resp);

  });
  useEffect(() => {
    const f = async () => {
      const resp = await client
        .get(`${baseUrl}/api/post/page/${page}`)
        .catch((e) => {
          console.log(e);
        });
      resp ? setPosts(resp.data as Post[]) : console.log(resp);
    };
    f();
  }, [page, location]);
  return (
    <form>
      <table>
        <thead>
          <tr>
            <th>
              削<br />除
            </th>
            <th>提供日</th>
            <th>UP職人</th>
            <th>提供品</th>
            <th>種別</th>
            <th>DL数</th>
            <th>トピック状況</th>
          </tr>
        </thead>
        <tbody>
          {posts.map((post) => (
            <tr key={`post-${post.id}`}>
              <td>
                <input type="checkbox" name={`${post.id}`} value="1" onChange={(e)=>{
                  if (e.target.checked){
                  setDeletePostIds([...delete_post_ids, post.id])
                } else{
                  setDeletePostIds( delete_post_ids.filter((v)=> v !== post.id))
                }
                }}/>
              </td>
              <td>{parseDate(post.uploaded_at!)}</td>
              <td>{post.user}</td>
              <td>
                <Link to={`/post/${post.id}`}>{post.title}</Link>
              </td>
              <td>{post.kind}</td>
              <td>{post.count ? post.count : 0}</td>
              <td>
                {generateTopicString({count: post.count,limit_count: post.limit_count,limit_date: post.limit_date,require_message: post.require_message})}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
      <div>
        <input type="password" onChange={(e)=>{
          setDeletePassword(e.target.value)
        }}></input>
        <input type="submit" value="削除" onClick={deletePost} />
        
      </div>
    </form>
  );
};

export default Content;
