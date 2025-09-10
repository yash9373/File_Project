import React, { useState, useEffect } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { useNavigate } from 'react-router-dom';
import { getMe, logout } from './redux/slice/authslice'; // Assuming you have this
import { fetchFiles, deleteFile, downloadFile } from './redux/slice/fileslice'; // Assuming you have this
import { createShareLink } from './redux/slice/linkslice'; // import your slice actions

// --- Helper function for date formatting ---
const formatDate = (dateString) => {
    const options = { year: 'numeric', month: 'long', day: 'numeric' };
    return new Date(dateString).toLocaleDateString(undefined, options);
};


// --- Reusable Components ---

// Modal Component
const Modal = ({ children, onClose }) => (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50 backdrop-blur-sm" onClick={onClose}>
        <div className="bg-white rounded-2xl shadow-xl p-6 sm:p-8 w-full max-w-md mx-4" onClick={e => e.stopPropagation()}>
            {children}
        </div>
    </div>
);


// Navbar Component
const Navbar = ({ user }) => {
    const dispatch = useDispatch();
    const navigate = useNavigate();

    const handleLogout = () => {
        dispatch(logout());
        navigate('/login');
    };

    return (
        <nav className="sticky top-4 inset-x-4 mx-auto max-w-7xl z-50 bg-white/70 backdrop-blur-lg rounded-2xl shadow-sm border border-gray-200">
            <div className="flex items-center justify-between px-6 py-3">
                <div className="flex items-center gap-4">
                    <div className="relative">
                        <img
                            className="w-11 h-11 rounded-full object-cover border-2 border-white shadow-md profile-icon-animated"
                            src={"https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRZF-VbuWgaUNg_BcV9h9nQwAEzkytE3ubdDg&s++"}
                            alt="User Avatar"
                        />
                    </div>
                    <span className="font-semibold text-gray-700 hidden sm:block">
                        {user?.name || "Loading..."}
                    </span>
                </div>
                <div className="flex items-center gap-6">
                    <a href="#" className="font-bold text-xl text-gray-800 tracking-tight">
                        <span>Secure</span><span className="text-blue-600">Vault</span>
                    </a>
                    <button
                        className="bg-red-500 text-white font-semibold px-5 py-2 rounded-lg hover:bg-red-600 transition-colors shadow"
                        onClick={handleLogout}
                    >
                        Log Out
                    </button>
                </div>
            </div>
        </nav>
    );
};


// "Add New" Card Component
const AddNewCard = ({ onClick }) => (
    <div onClick={onClick} className="flex items-center justify-center p-6 bg-slate-200 border-2 border-dashed border-slate-400 rounded-xl cursor-pointer hover:bg-slate-300 hover:border-slate-500 transition-colors">
        <div className="text-center text-slate-600">
            <svg xmlns="http://www.w3.org/2000/svg" className="mx-auto h-12 w-12" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth="1.5">
                <path strokeLinecap="round" strokeLinejoin="round" d="M12 4v16m8-8H4" />
            </svg>
            <p className="mt-2 font-semibold">Upload File</p>
        </div>
    </div>
);

// Standard Content Card Component
const ContentCard = ({ file, onDownload, onDelete, onShare }) => (
    <div className="bg-white p-5 rounded-xl shadow-sm border border-gray-200 flex flex-col justify-between">
        <div>
            <div className="w-full h-24 bg-gray-200 rounded-lg mb-4 flex items-center justify-center">
                 <svg xmlns="http://www.w3.org/2000/svg" className="h-10 w-10 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth="1">
                    <path strokeLinecap="round" strokeLinejoin="round" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                </svg>
            </div>
            <h3 className="font-semibold text-gray-800 truncate">{file.filename}</h3>
            <p className="text-sm text-gray-500">Added: {formatDate(file.created_at)}</p>
        </div>
        <div className="mt-4 flex gap-2">
             <button onClick={() => onDownload(file.id)} className="flex-1 text-sm bg-blue-500 text-white font-semibold px-3 py-1.5 rounded-md hover:bg-blue-600 transition-colors shadow-sm">⬇</button>
             <button onClick={() => onShare(file.id)} className="text-sm bg-gray-200 text-gray-700 font-semibold px-3 py-1.5 rounded-md hover:bg-gray-300 transition-colors">📤
</button>
             <button onClick={() => onDelete(file.id)} className="bg-red-100 text-red-600 p-1.5 rounded-md hover:bg-red-200 transition-colors">
                <svg xmlns="http://www.w3.org/2000/svg" className="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth="2">
                    <path strokeLinecap="round" strokeLinejoin="round" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                </svg>
            </button>
        </div>
    </div>
);


// --- Main App Component ---
export default function Layout() {
    const dispatch = useDispatch();
    const navigate = useNavigate();

    // --- State from Redux ---
    const { user, token: authToken } = useSelector(state => state.auth);
    const { files, loading, error } = useSelector(state => state.file);

    const [isDownloadModalOpen, setDownloadModalOpen] = useState(false);
    const [isShareModalOpen, setShareModalOpen] = useState(false);
    const [selectedFileId, setSelectedFileId] = useState(null);
    const [password, setPassword] = useState('');
    const [shareLink, setShareLink] = useState('');


    useEffect(() => {
        if (!authToken) {
            navigate('/login');
        } else {
            dispatch(getMe());
            dispatch(fetchFiles());
        }
    }, [authToken, dispatch, navigate]);

    // --- Event Handlers ---
    const handleDownloadClick = (fileId) => {
        setSelectedFileId(fileId);
        setDownloadModalOpen(true);
        setPassword('');
    };

    const handleDelete = (fileId) => {
        if (window.confirm("Are you sure you want to delete this file?")) {
            dispatch(deleteFile(fileId));
        }
    };const handleShareClick = async (fileId) => {
    setSelectedFileId(fileId);
    setShareModalOpen(true);
    setShareLink("Generating link...");

    try {
        const resultAction = await dispatch(createShareLink({ file_id: fileId }));
        
        if (createShareLink.fulfilled.match(resultAction)) {
            // The action payload already contains the full, correct link
            setShareLink(resultAction.payload.link);
        } else {
            const errorMessage = resultAction.payload?.error || "Failed to generate link.";
            setShareLink(errorMessage);
        }
    } catch (err) {
        console.error(err);
        setShareLink("Error generating link.");
    }
};

    const handleConfirmDownload = async () => {
        if (!password) {
            alert("Password is required.");
            return;
        }
        try {
            const blob = await downloadFile(selectedFileId, password, authToken);
            const file = files.find(f => f.id === selectedFileId);
            const url = window.URL.createObjectURL(blob);
            const link = document.createElement('a');
            link.href = url;
            link.setAttribute('download', file?.filename || 'download');
            document.body.appendChild(link);
            link.click();
            link.parentNode.removeChild(link);
            setDownloadModalOpen(false);
        } catch (err) {
            alert("Download failed. Check your password or try again.");
            console.error(err);
        }
    };


    return (
        <div className="bg-slate-100 min-h-screen">
            <style>{`
                @keyframes pulse { 0%, 100% { box-shadow: 0 0 0 0 rgba(37, 99, 235, 0.4); } 70% { box-shadow: 0 0 0 10px rgba(37, 99, 235, 0); } }
                .profile-icon-animated { animation: pulse 2s infinite; }
                body::-webkit-scrollbar { width: 8px; }
                body::-webkit-scrollbar-track { background: #f1f5f9; }
                body::-webkit-scrollbar-thumb { background: #cbd5e1; border-radius: 10px; }
                body::-webkit-scrollbar-thumb:hover { background: #94a3b8; }
            `}</style>

            <Navbar user={user} />

            <main className="pt-8 pb-16 px-4 sm:px-6 lg:px-8">
                <div className="max-w-7xl mx-auto">
                    {loading && <p className="text-center">Loading files...</p>}
                    {error && <p className="text-center text-red-500">{error}</p>}

                    <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-6">
                        <AddNewCard onClick={() => navigate('/upload')} />
                        {files.map(file => (
                            <ContentCard
                                key={file.id}
                                file={file}
                                onDownload={handleDownloadClick}
                                onDelete={handleDelete}
                                onShare={handleShareClick}
                            />
                        ))}
                    </div>
                </div>
            </main>

            {/* Download Modal */}
            {isDownloadModalOpen && (
                <Modal onClose={() => setDownloadModalOpen(false)}>
                    <h3 className="text-lg font-bold text-gray-800 mb-4">Enter Password to Download</h3>
                    <p className="text-sm text-gray-600 mb-4">This file is encrypted. Please enter the password you used during upload.</p>
                    <input
                        type="password"
                        value={password}
                        onChange={(e) => setPassword(e.target.value)}
                        placeholder="File password"
                        className="w-full px-4 py-2 border border-gray-300 rounded-lg mb-4 focus:ring-2 focus:ring-blue-500 text-black"
                    />
                    <div className="flex justify-end gap-3">
                        <button onClick={() => setDownloadModalOpen(false)} className="px-4 py-2 rounded-lg bg-gray-200 text-gray-700 font-semibold hover:bg-gray-300">Cancel</button>
                        <button onClick={handleConfirmDownload} className="px-4 py-2 rounded-lg bg-blue-600 text-white font-semibold hover:bg-blue-700">Confirm</button>
                    </div>
                </Modal>
            )}

            {/* Share Modal */}
            {isShareModalOpen && (
                 <Modal onClose={() => setShareModalOpen(false)}>
                    <h3 className="text-lg font-bold text-gray-800 mb-4">Share File</h3>
                    <p className="text-sm text-gray-600 mb-4">Here is your shareable link. Remember, the recipient will still need the file password.</p>
                     <input
                        type="text"
                        readOnly
                        value={shareLink}
                        className="w-full px-4 py-2 border bg-gray-100 border-gray-300 rounded-lg mb-4 text-black"
                    />
                    <div className="flex justify-end">
                         <button onClick={() => setShareModalOpen(false)} className="px-4 py-2 rounded-lg bg-gray-200 text-gray-700 font-semibold hover:bg-gray-300">Close</button>
                    </div>
                </Modal>
            )}
        </div>
    );
}


