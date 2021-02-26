/**
 * Represents a user from the database table.
 */
interface User {
  id?: number;
  name?: string;
  email?: string;
  pass?: string;
  created_on?: Date;
}

export { User };
